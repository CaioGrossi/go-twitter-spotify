package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)


func connectTwitterHttpStream(client http.Client) {
	req , err := http.NewRequest("GET", "https://api.twitter.com/2/tweets/search/stream", nil)

	if err != nil {
		fmt.Printf("Error in start twitter stream: %s\n", err)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + os.Getenv("TWITTER_TOKEN")},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error in start twitter stream: %s\n", err)
	}

	defer resp.Body.Close()

	type Tweet struct {
		Id, Text string
	}

	type Message struct {
		Tweet Tweet
	}

	tweet := new(Message);

	decoder := json.NewDecoder(resp.Body)

	for {
		err := decoder.Decode(&tweet)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(tweet.Tweet.Text)

		addTrackToPlaylist(client, tweet.Tweet.Text)
	}
}
