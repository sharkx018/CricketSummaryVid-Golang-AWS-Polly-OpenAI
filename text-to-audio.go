package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"log"
	"os"
)

func textToAudio(text string, id string) {

	fileName := fmt.Sprintf("audio-%s.mp3", id)

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

	// Create a Polly client
	svc := polly.New(sess)

	// Text to convert to speech
	//***text := "man playing cricket match in stadium"

	// Specify the voice and output format
	voiceID := "Joanna"
	outputFormat := "mp3"

	// Generate the request to synthesize speech
	input := &polly.SynthesizeSpeechInput{
		Text:         aws.String(text),
		OutputFormat: aws.String(outputFormat),
		VoiceId:      aws.String(voiceID),
	}

	// Send the request to Polly
	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		log.Fatalf("Failed to convert text to audio: %v", err)
	}

	// Create an output file to save the audio
	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Copy the audio data to the output file
	_, err = io.Copy(outputFile, output.AudioStream)
	if err != nil {
		log.Fatalf("Failed to save audio data to file: %v", err)
	}

	fmt.Println("Text converted to audio and saved as " + fileName)

}
