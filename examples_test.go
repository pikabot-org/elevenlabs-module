package elevenlabs_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/pikabot-org/elevenlabs-module"
)

var apiKey string = "YOUR_API_KEY_HERE" // todo: Find a way to get the API key from the environment

func TestExampleClient_TextToSpeech(t *testing.T) {
	// Create a new client
	client := elevenlabs.NewClient(context.Background(), apiKey, 30*time.Second)

	// Create a TextToSpeechRequest
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    "Hello, world! My name is Adam, nice to meet you!",
		ModelID: "eleven_monolingual_v1",
	}

	// Call the TextToSpeech method on the client, using the "Adam"'s voice ID.
	audio, err := client.TextToSpeech("pNInz6obpgDQGcFmaJgB", ttsReq)
	if err != nil {
		t.Fatalf("Got %T error: %q\n", err, err)
		return
	}

	// Write the audio file bytes to disk
	if err := os.WriteFile("adam.mp3", audio, 0644); err != nil {
		t.Fatalf("Could not write audio file: %s", err)
		return
	}

	t.Log("Successfully generated audio file")
}

func TestExampleClient_TextToSpeechStream(t *testing.T) {
	message := `The concept of "flushing" typically applies to I/O buffers in many programming 
languages, which store data temporarily in memory before writing it to a more permanent location
like a file or a network connection. Flushing the buffer means writing all the buffered data
immediately, even if the buffer isn't full.`

	// Set your API key
	elevenlabs.SetAPIKey(apiKey)

	// Set a large enough timeout to ensure the stream is not interrupted.
	elevenlabs.SetTimeout(1 * time.Minute)

	// We'll use mpv to play the audio from the stream piped to standard input
	cmd := exec.CommandContext(context.Background(), "mpv", "--no-cache", "--no-terminal", "--", "fd://0")

	// Get a pipe connected to the mpv's standard input
	pipe, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Could not get pipe: %s", err)
		return
	}

	// Attempt to run the command in a separate process
	if err := cmd.Start(); err != nil {
		t.Fatalf("Could not start mpv: %s", err)
	}

	// Stream the audio to the pipe connected to mpv's standard input
	if err := elevenlabs.TextToSpeechStream(
		pipe,
		"pNInz6obpgDQGcFmaJgB",
		elevenlabs.TextToSpeechRequest{
			Text:    message,
			ModelID: "eleven_multilingual_v1",
		}); err != nil {
		t.Fatalf("Got %T error: %q\n", err, err)
	}

	// Close the pipe when all stream has been copied to the pipe
	if err := pipe.Close(); err != nil {
		t.Fatalf("Could not close pipe: %s", err)
	}
	t.Log("Streaming finished.")

	// Wait for mpv to exit. With the pipe closed, it will do that as
	// soon as it finishes playing
	if err := cmd.Wait(); err != nil {
		t.Fatalf("mpv exited with error: %s", err)
		return
	}

	t.Log("All done.")
}

func TestExampleClient_GetHistory(t *testing.T) {
	// Define a helper function to print history items
	printHistory := func(r elevenlabs.GetHistoryResponse, p int) {
		fmt.Printf("--Page %d--\n", p)
		for i, h := range r.History {
			t := time.Unix(int64(h.DateUnix), 0)
			fmt.Printf("%d. %s - %s: %d bytes\n", p+i, t.Format("2006-01-02 15:04:05"), h.HistoryItemId, len(h.Text))
		}
	}
	// Create a new client
	client := elevenlabs.NewClient(context.Background(), apiKey, 30*time.Second)

	// Get and print the first page (5 items).
	page := 1
	historyResp, nextPage, err := client.GetHistory(elevenlabs.PageSize(5))
	if err != nil {
		t.Fatalf("Got %T error: %q\n", err, err)
	}
	printHistory(historyResp, page)

	// Get all other pages
	for nextPage != nil {
		page++
		// Retrieve the next page. The page size from the original call is retained but
		// can be overwritten by passing a call to PageSize with the new size.
		historyResp, nextPage, err = nextPage()
		if err != nil {
			t.Fatalf("Got %T error: %q\n", err, err)
		}
		printHistory(historyResp, page)
	}
}

func TestExampleClient_SpeechToSpeech(t *testing.T) {
	// Create a new client
	client := elevenlabs.NewClient(context.Background(), apiKey, 30*time.Second)
	// Using previously generated audio file "adam.mp3" as input
	inputAudio, err := os.Open("adam.mp3")
	if err != nil {
		t.Fatalf("Could not open audio file: %s", err)
		return
	}
	defer inputAudio.Close()

	// Create a SpeechToSpeechRequest
	stsReq := elevenlabs.SpeechToSpeechRequest{
		Audio:   inputAudio,
		ModelID: "eleven_english_sts_v2",
	}

	// Call the SpeechToSpeech method on the client, using the "Rachel"'s voice ID.
	audio, err := client.SpeechToSpeech("21m00Tcm4TlvDq8ikWAM", stsReq)
	if err != nil {
		t.Fatalf("Got %T error: %q\n", err, err)
		return
	}

	// Write the audio file bytes to disk
	if err := os.WriteFile("rachel.mp3", audio, 0644); err != nil {
		t.Fatalf("Could not write audio file: %s", err)
		return
	}

	t.Log("Successfully generated audio file")
}
