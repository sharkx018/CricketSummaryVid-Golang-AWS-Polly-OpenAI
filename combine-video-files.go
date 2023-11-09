package main

import (
	"fmt"
	"os"
	"os/exec"
)

// /Users/mukulverma/Downloads/ffmpeg -i o1.mp4 -i o2.mp4 -i o3.mp4 -i o4.mp4 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a][3:v][3:a]concat=n=4:v=1:a=1[vout][aout]" -map "[vout]" -map "[aout]" -strict -2 output.mp4

func CombineChunks(inputFiles []string) {

	// cmd := exec.Command("/Users/mukulverma/Downloads/ffmpeg", "-i", "o1.mp4", "-i", "o2.mp4", "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[vout][aout]", "-map", "[vout]", "-map", "[aout]", "-strict", "-2", "output.mp4")
	// shell cmd   /Users/mukulverma/Downloads/ffmpeg -i o1.mp4 -i o2.mp4 -i o3.mp4 -i o4.mp4 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a][3:v][3:a]concat=n=4:v=1:a=1[vout][aout]" -map "[vout]" -map "[aout]" -strict -2 output.mp4

	// Input video file paths
	//inputFiles := []string{"o1.mp4", "o2.mp4", "o3.mp4"} // Add more files as needed

	args := []string{}

	fileName := "finalVideo.mp4"

	for _, inputFile := range inputFiles {
		args = append(args, "-i", inputFile)
	}

	//for _, inputFile := range inputFiles {
	//	args = append(args, "-i", "tmp/"+inputFile)
	//}

	args = append(args, "-filter_complex")

	tmpStr := ""
	for i, _ := range inputFiles {
		tmpStr += fmt.Sprintf("[%d:v][%d:a]", i, i)
	}
	tmpStr += fmt.Sprintf("concat=n=%d:v=1:a=1[vout][aout]", len(inputFiles))

	args = append(args, tmpStr)

	args = append(args, "-map", "[vout]", "-map", "[aout]", "-strict", "-2", fileName)

	//cmd := exec.Command("/Users/mukulverma/Downloads/ffmpeg", args...)
	cmd := exec.Command(FFMPEG_PATH, args...)
	// Set the output and error streams to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the FFmpeg command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
