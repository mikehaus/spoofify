package helpers

// Util file to handle oauth2 initialization and config for Spotify

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	SpotifyAuthUrl  = "https://accounts.spotify.com/authorize?"
	SpotifyTokenUrl = "https://accounts.spotify.com/api/token"
	RedirectUri     = "localhost:8080/auth/spotify/callback"
)

const (
	ScopeStreaming           = "streaming"
	ScopeReadPlaylistPrivate = "playlist-read-private"
	ScopePlaybackPosition    = "user-read-playback-position"
)

type SpotifyAuth struct {
	Config *oauth2.Config
	Token  *oauth2.Token
}

func NewSpotifyAuth() *SpotifyAuth {
	return &SpotifyAuth{
		Config: SpotifyOAuthConfig(),
	}
}

func SpotifyOAuthConfig() *oauth2.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

func (auth *SpotifyAuth) SpotifyOAuthUrl() string {
	verifier := oauth2.GenerateVerifier()
	return auth.Config.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
}

func (auth *SpotifyAuth) SpotifyClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return auth.Config.Client(ctx, token)
}

func (auth *SpotifyAuth) SpotifyAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	code := r.URL.Query().Get("code")

	token, err := auth.Config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	auth.Token = token

  // client := auth.Config.Client(ctx, token)
  // TODO: Do stuff with client below
}
