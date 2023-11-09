package main

import (
	"fmt"
	"strconv"
)

//const (
//	AWS_ACCES_ID = "****"
//	AWS_SECRET   = "****"
//	REGION       = "****"
//
//	OPEN_API_KEY = "****"
//	OPEN_API_URL = "https://api.openai.com/v1/images/generations"
//	SIZE         = "1024x1024"
//)

const (
	FFMPEG_PATH = "/Users/mukulverma/Downloads/ffmpeg"
)

func main() {

	//commentary := []string{
	//	"a boy is playing",
	//	"a dog is barking",
	//}

	// get data from scraper
	//matchLink := "https://www.cricbuzz.com/cricket-scores/82374/nep-vs-sin-1st-match-group-a-icc-mens-t20i-world-cup-asia-finals-2023"
	matchLink := "https://www.cricbuzz.com/cricket-scores/82376/uae-vs-bhr-2nd-match-group-b-icc-mens-t20i-world-cup-asia-finals-2023"

	commentary, err := scrapeData(matchLink)
	if err != nil {
		fmt.Println("Error from scrapeData:", err.Error())
		return
	}

	chunks := []string{}

	for id, text := range commentary {

		// create the audio file
		textToAudio(text, strconv.Itoa(id))

		// create the image file
		downloadImage(text, strconv.Itoa(id))

		// create the chunk file
		createVideoChunk(strconv.Itoa(id))

		// append to chunk array
		chunks = append(chunks, fmt.Sprintf("chunk-%d.mp4", id))
	}

	// combine all chunks to video
	CombineChunks(chunks)

	//delete temporary files
	//DeleteTempFiles(commentary)

}
