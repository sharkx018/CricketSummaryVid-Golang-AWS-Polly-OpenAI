package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func scrapeData() {
	url := "https://www.cricbuzz.com/cricket-match-highlights/75595/eng-vs-ned-40th-match-icc-cricket-world-cup-2023"

	commentary, err := scrapeCommentary(url)
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range commentary {
		fmt.Println(comment)
	}
}

func scrapeCommentary(url string) ([]string, error) {
	var commentary []string

	// Make an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return commentary, err
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return commentary, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the HTML content
	htmlContent, err := readResponseBody(resp)
	if err != nil {
		return commentary, err
	}

	// Define a regular expression to match the desired format
	re := regexp.MustCompile(`^\d+\.\d+\s[\w\s]+, .+`)

	// Find and store text that matches the format
	matches := re.FindAllString(htmlContent, -1)
	commentary = append(commentary, matches...)

	return commentary, nil
}

func readResponseBody(resp *http.Response) (string, error) {
	buf := make([]byte, 1024)
	var body string

	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return "", err
		}
		body += string(buf[:n])
	}

	return body, nil
}
