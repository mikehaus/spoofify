package helpers

import (
  "time"
	"math/rand"
	"golang.org/x/oauth2"
)

const (
	SpotifyAuthUrl = "https://accounts.spotify.com/authorize?"
	RedirectUri    = "http://localhost:8080/callback"
)

const (
  ScopeStreaming = "streaming"
  ScopeReadPlaylistPrivate = "playlist-read-private"
  ScopePlaybackPosition = "user-read-playback-position"
)

func GetAuthUrl() string {
  return SpotifyAuthUrl
}

// TODO: May need to initialize this on the top level and move this up
func initSeed() {
  rand.Seed(time.Now().UnixNano())
}

func GetUriAuthState() string {
  initSeed()
	return randSeq(16)
}

// Helper used to generate random state value for Spotify auth req
func randSeq(n int) string {
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
