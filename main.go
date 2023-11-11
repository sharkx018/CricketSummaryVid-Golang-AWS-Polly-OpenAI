package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"strconv"
)

// create your api-keys
//const (
//	AWS_ACCES_ID = "****"
//	AWS_SECRET   = "****"
//	REGION       = "****"
//
//	OPEN_API_KEY = "****"
//	OPEN_API_URL = "https://api.openai.com/v1/images/generations"
//	SIZE         = "1024x1024"
//)

var awsSess *session.Session

const (
	FfmpegPath = "/Users/mukulverma/Downloads/ffmpeg"
)

func init() {
	initAWSSession()
}

func main() {

	// get data from scraper
	matchLink := "https://www.cricbuzz.com/cricket-scores/82376/uae-vs-bhr-2nd-match-group-b-icc-mens-t20i-world-cup-asia-finals-2023"

	commentary, err := scrapeData(matchLink)
	if err != nil {
		fmt.Println("Error from scrapeData:", err.Error())
		return
	}

	var chunks []string

	for id, text := range commentary {

		// create the audio file
		textToAudio(text, strconv.Itoa(id), awsSess)

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

func initAWSSession() {

	// Specify the AWS region and create a session
	region := REGION // Change this to your desired AWS region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(AWS_ACCES_ID,
			AWS_SECRET, ""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	awsSess = sess

}
