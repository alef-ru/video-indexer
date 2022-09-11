package embds

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const cmd = "ffmpeg"

type FrameInfo struct {
	//the is no frame number field, it's supposed to be the index of []FrameInfo
	Time       time.Duration
	SceneScore float64
}

type FfmpegWrapper struct {
	execWrapper execWrapper
}

func NewFfmpegWrapper() (*FfmpegWrapper, error) {
	ew, err := newExecWrapper(cmd)
	if err != nil {
		return nil, err
	}
	ew.enableStderrForwarding()
	return &FfmpegWrapper{execWrapper: *ew}, nil
}

func (fw *FfmpegWrapper) FrameDiffs(ctx context.Context, videoPathOrUrl string) ([]FrameInfo, error) {
	tmpFile, err := os.CreateTemp("", "frame_diffs_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("ERROR Failed to delete temp file %s, %v", name, err)
		}
	}(tmpFile.Name())
	log.Printf("INFO temp file created: %s", tmpFile.Name())
	err = fw.execWrapper.execute(ctx, "-i", videoPathOrUrl, "-vf",
		fmt.Sprintf("select='gte(scene,0)',metadata=print:file=%s", tmpFile.Name()),
		"-an", "-f", "null", "-")
	if err != nil {
		return nil, fmt.Errorf("failed to execute ffmpeg: %w", err)
	}
	return fw.parseFrameDiffs(tmpFile)
}

func (fw *FfmpegWrapper) parseFrameDiffs(file *os.File) ([]FrameInfo, error) {
	var (
		frames             []FrameInfo
		lineNum            int
		curFrameTimestamp  time.Duration
		curFrameSceneScore float64
		err                error //to avoid shadowing
	)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if lineNum%2 == 0 { //even lines look like `frame:27   pts:81081   pts_time:0.9009`
			curFrameTimestamp, err = parseFrameTimestamp(line, lineNum)
		} else { // odd lines look like `lavfi.scene_score=0.000014`
			curFrameSceneScore, err = parseSceneScore(line, lineNum)
			frames = append(frames, FrameInfo{curFrameTimestamp, curFrameSceneScore})
		}
		if err != nil {
			return nil, err
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return frames, nil
}

var frameDiffEvenLineRE = regexp.MustCompile(`^frame:(\d+)\s+pts:(\d+)\s+pts_time:(\d+(\.\d+)?)$`)

func parseFrameTimestamp(line string, lineNum int) (time.Duration, error) {
	params := frameDiffEvenLineRE.FindStringSubmatch(line)
	if len(params) != 5 {
		return 0, fmt.Errorf("failed to parse line %d: %s", lineNum, line)
	}
	frameNum, err := strconv.Atoi(params[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse frame number. Line %d: %s", lineNum, line)
	}
	if lineNum/2 != frameNum {
		return 0, fmt.Errorf("unexpected frame number. Line %d: %s", lineNum, line)
	}
	ptsTime, err := strconv.ParseFloat(params[3], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse pts_time. Line %d: %s", lineNum, line)
	}
	return time.Duration(ptsTime * float64(time.Second)), nil
}

var frameDiffOddLineRE = regexp.MustCompile(`^lavfi.scene_score=(\d+(\.\d+)?)$`)

func parseSceneScore(line string, lineNum int) (float64, error) {
	params := frameDiffOddLineRE.FindStringSubmatch(line)
	if len(params) != 3 {
		return 0, fmt.Errorf("failed to parse line %d: %s", lineNum, line)
	}
	return strconv.ParseFloat(params[1], 64)
}

func (ffmpeg *FfmpegWrapper) EnableStderrForwarding() {
	ffmpeg.execWrapper.enableStderrForwarding()
}

func (ffmpeg *FfmpegWrapper) DisableStderrForwarding() {
	ffmpeg.execWrapper.disableStderrForwarding()
}
