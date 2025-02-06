package helpers

// Go server to handle REST operations for Spotify API

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

func InitServer(auth *SpotifyAuth) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: serve(auth),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func serve(auth *SpotifyAuth) http.Handler {
	mux := http.NewServeMux()

	// oauth Spotify handlers
	mux.HandleFunc("/auth/spotify/login", auth.HandleSpotifyLogin)
	mux.HandleFunc("/auth/spotify/callback", auth.SpotifyAuthCallback)

	return mux
}

// TODO: may not need cookie?
func generateStateOAuthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	bitStr := make([]byte, 16)
	rand.Read(bitStr)
	state := base64.URLEncoding.EncodeToString(bitStr)
	cookie := http.Cookie{Name: "spotifyoauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
