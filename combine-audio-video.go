package main

import (
	"fmt"
	"os/exec"
)

func createVideoChunk(id string) {

	// Input audio file (replace with the path to your audio file)
	audioFile := "audio-" + id + ".mp3"

	// Input image file (replace with the path to your image file)
	imageFile := "image-" + id + ".png"

	// Output video file
	outputVideoFile := "chunk-" + id + ".mp4" // Change the format if needed

	// Full path to the ffmpeg executable (replace with the actual path)
	//ffmpegPath := "/Users/mukulverma/Downloads/ffmpeg"
	ffmpegPath := FfmpegPath

	// Use FFmpeg to combine the audio and image into a video
	cmd := exec.Command(ffmpegPath,
		"-loop", "1",
		"-i", imageFile,
		"-i", audioFile,
		"-c:v", "libx264",
		"-t", "5",
		"-pix_fmt", "yuv420p",
		"-vf", "scale=1920:1080",
		"-c:a", "aac",
		"-strict", "experimental",
		"-shortest",
		outputVideoFile,
	)

	// Run the FFmpeg command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running FFmpeg:", err)
		return
	}

	fmt.Printf("Video created: %s\n", outputVideoFile)
}
