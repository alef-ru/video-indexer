package embds_test

import (
	"context"
	"fmt"
	"github.com/alef-ru/vindex/embds"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestFfmpegWrapper_FrameDiffs(t *testing.T) {
	ffmpeg, err := embds.NewFfmpegWrapper()
	if err != nil {
		t.Fatalf("failed to init ffmpeg wrapper: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	frames, err := ffmpeg.FrameDiffs(ctx, "testdata/countdown.mp4")
	assert.NoError(t, err)
	assert.Equal(t, 100, len(frames))
	for i, frame := range frames {
		if frame.Time != 0 && math.Abs(math.Round(frame.Time.Seconds())-frame.Time.Seconds()) < 0.02 {
			// Picture changes every second
			fmt.Printf("frame %d: %+v\n", i, frame)
			assert.Greater(t, frame.SceneScore, 0.5)
		} else {
			// Other frames have relatively low scene change
			assert.Less(t, frame.SceneScore, 0.5)
		}

	}
}
