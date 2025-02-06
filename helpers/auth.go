package helpers

// Util file to handle oauth2 initialization and config for Spotify

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const (
	// SpotifyAuthUrl  = "https://accounts.spotify.com/authorize?"
	// SpotifyTokenUrl = "https://accounts.spotify.com/api/token"
	RedirectUri = "localhost:8080/auth/spotify/callback"
)

const (
	ScopeStreaming           = "streaming"
	ScopeReadPlaylistPrivate = "playlist-read-private"
	ScopePlaybackPosition    = "user-read-playback-position"
)

type SpotifyAuth struct {
	config *oauth2.Config
	state  string
	token  *oauth2.Token
}

func NewSpotifyAuth() *SpotifyAuth {
	return &SpotifyAuth{
		config: SpotifyOAuthConfig(),
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
		// Endpoint:     oauth2.Endpoint{AuthURL: SpotifyAuthUrl, TokenURL: SpotifyTokenUrl},
		Endpoint: spotify.Endpoint,
	}

	return config
}

func (auth *SpotifyAuth) SpotifyOAuthUrl() string {
	verifier := oauth2.GenerateVerifier()
	return auth.config.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
}

func (auth *SpotifyAuth) SpotifyClient(ctx context.Context) *http.Client {
	return auth.config.Client(ctx, auth.token)
}

func (auth *SpotifyAuth) HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.config.AuthCodeURL(auth.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (auth *SpotifyAuth) SpotifyAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	fmt.Printf("Found the code: %s", code)

	if state != auth.state {
		http.Error(w, "Invalid State", http.StatusBadRequest)
		return
	}

	token, err := auth.config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		log.Println("Token exchange error:", err)
		return
	}

	auth.token = token
}

func generateRandomState() string {
	return fmt.Sprintf("%d", rand.Intn(100000))
}
