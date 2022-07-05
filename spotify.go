package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func authenticateSpotifyApi(client http.Client) (string) {
	auth_data := url.Values{}
    auth_data.Set("grant_type", "refresh_token")
    auth_data.Set("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	req , err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(auth_data.Encode()))

	req.Header = http.Header{
		"Authorization": {"Basic " + os.Getenv("SPOTIFY_INITIAL_TOKEN")},
		"Content-Type": {"application/x-www-form-urlencoded"},
	}

	if err != nil {
		fmt.Printf("Error getting spotify token: %s\n", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error getting spotify token: %s\n", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	type AuthReturn struct {
		AccessToken string `json:"access_token"`
	}

	authReturn := AuthReturn{}

	jsonErr := json.Unmarshal(body, &authReturn)

	if jsonErr != nil {
		fmt.Println(err)
	}

	return authReturn.AccessToken
}

type SearchReturn struct {
	Tracks Track `json:"tracks"`
}

type Track struct {
	Items []Item  `json:"items"`
}

type Item struct {
	Uri string `json:"uri"`
}


func searchTrackAtSpotifyApi(client http.Client, track string, authToken string) string {
	req , err := http.NewRequest("GET", "https://api.spotify.com/v1/search?type=track&q=" + strings.Replace(track, " ", "+", -1), nil)

	if err != nil {
		fmt.Printf("Error getting track: %s\n", err)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + authToken},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error getting spotify token: %s\n", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	searchReturn := SearchReturn{}

	jsonErr := json.Unmarshal(body, &searchReturn)

	if jsonErr != nil {
		fmt.Println(err)
	}

	if len(searchReturn.Tracks.Items) > 0 {
		return searchReturn.Tracks.Items[0].Uri
	}

	return ""
}

func addTrackToPlaylist(client http.Client, track string) bool {

	authToken := authenticateSpotifyApi(client)

	trackUri := searchTrackAtSpotifyApi(client, track, authToken)

	if trackUri == "" {
		return false
	}

	var trackData = []byte(`{
		"uris": [
			"` + trackUri + `"
		]
	}`)

	req, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/60VbGsXUSdOTM2txWPFiFe/tracks", bytes.NewBuffer(trackData))

	if err != nil {
		fmt.Printf("Error getting track: %s\n", err)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + authToken},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error getting spotify token: %s\n", err)
	}

	if resp.StatusCode != 201 {
		return false
	}

	return true
}
