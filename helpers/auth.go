package helpers

// Util file to handle oauth2 initialization and config for Spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const (
	SpotifyAuthUrl  = "https://accounts.spotify.com/authorize?"
	SpotifyTokenUrl = "https://accounts.spotify.com/api/token"
	RedirectUri     = "http://localhost:8080/callback"
)

const (
	ScopeStreaming           = "streaming"
	ScopeReadPlaylistPrivate = "playlist-read-private"
	ScopePlaybackPosition    = "user-read-playback-position"
)

// TODO: handle client and url separately
func GetSpotifyOAuthClient() (*http.Client, string) {
	ctx := context.Background()

	scopes := []string{
		ScopeStreaming,
		ScopePlaybackPosition,
		ScopeReadPlaylistPrivate,
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Scopes:       scopes,
		Endpoint:     oauth2.Endpoint{AuthURL: SpotifyAuthUrl, TokenURL: SpotifyTokenUrl},
	}

	verifier := oauth2.GenerateVerifier()
	url := config.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	return config.Client(ctx, token), url
}
