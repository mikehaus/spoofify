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
	RedirectUri     = "http://localhost:8080/auth/spotify/callback"
)

const (
	ScopeStreaming           = "streaming"
	ScopeReadPlaylistPrivate = "playlist-read-private"
	ScopePlaybackPosition    = "user-read-playback-position"
)

type SpotifyAuth struct {
	OauthConfig *oauth2.Config
}

func NewSpotifyAuth() *SpotifyAuth {
	return &SpotifyAuth{
		OauthConfig: SpotifyOAuthConfig(),
	}
}

func SpotifyOAuthConfig() *oauth2.Config {
	scopes := []string{
		ScopeStreaming,
		ScopePlaybackPosition,
		ScopeReadPlaylistPrivate,
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  RedirectUri,
		Scopes:       scopes,
		Endpoint:     oauth2.Endpoint{AuthURL: SpotifyAuthUrl, TokenURL: SpotifyTokenUrl},
	}

	return config
}

func SpotifyOAuthUrl(config *oauth2.Config) string {
	verifier := oauth2.GenerateVerifier()
	return config.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
}

func SpotifyClient(config oauth2.Config) *http.Client {
	ctx := context.Background()

	var code string
	fmt.Scan(&code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	return config.Client(ctx, token)
}
