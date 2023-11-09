package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Item struct {
	RevisedPrompt string `json:"revised_prompt"`
	Url           string `json:"url"`
}

type Response struct {
	Data []Item `json:"data"`
}

func downloadImage(prompt, id string) {

	// API endpoint URL
	apiURL := OPEN_API_URL

	// API key
	apiKey := OPEN_API_KEY

	model := "dall-e-3"
	// prompt := "man is playing football"
	n := 1
	size := SIZE

	// Create a map to represent the payload
	payload := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"n":      n,
		"size":   size,
	}

	// Convert the payload to a JSON string
	requestPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", resp.StatusCode)
		return
	}

	var apiResponse Response
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	downloadImageFile(apiResponse.Data[0].Url, id)

	// Handle the response here (e.g., save the image or process it)

	fmt.Println("Request was successful. Handle the response as needed.")
}

func downloadImageFile(imageURL string, id string) {
	// URL of the image you want to download
	//imageURL := "https://oaidalleapiprodscus.blob.core.windows.net/private/org-vKZFHD9eQ1fsYRjYAUzxwZYR/user-QFaopBC9yb1Jlq2jZvTDn53F/img-Xe6ku3H0rjcJz69bQN5vGgyr.png?st=2023-11-08T17%3A47%3A07Z&se=2023-11-08T19%3A47%3A07Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-11-08T18%3A20%3A03Z&ske=2023-11-09T18%3A20%3A03Z&sks=b&skv=2021-08-06&sig=jnHcVvLOlQ2ZGdfVczPpvT3zBv5ZA81aHfKn8YjqEu0%3D"

	// Specify the output file path and name
	outputFile := "image-" + id + ".png" // Replace with the desired local file name

	// Create an HTTP GET request to the image URL
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", resp.StatusCode)
		return
	}

	// Create a new file to save the image
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Copy the image data from the HTTP response to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error copying image data:", err)
		return
	}

	fmt.Printf("Image downloaded and saved as '%s'\n", outputFile)
}
