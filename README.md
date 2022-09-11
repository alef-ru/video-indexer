# vindex - Video Indexer

It's one more project I'm writing to gain some experience with Go.
When (or if) I finish it, it will be a service
that can build the index of video files, extracting some features from the video content.

Typical use case: you have an archive of videos and want to search for some video clips from social networks against this archive.

It's supposed to be tolerant to changes in the video, which is usually done when it is published on a social network
* bitrate changes
* framerate changes
* resolution changes
* adding watermarks
* cutting and concatenating videos to some extent

But the index is supposed to search for an exact content match, not similar videos.
So, if you crop, mirror video, put it inside another picture, or apply some filters to the video,
the index should not match the altered video to the original one  
