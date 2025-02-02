package helpers

// Go server to handle REST operations for Spotify API

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"mikehaus/spoofify/helpers"
	"net/http"
	"time"
)

func server() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: serve(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("%v", err)
	}
}

func serve() http.Handler {
	mux := http.NewServeMux()

	// oauth Spotify handlers
	mux.HandleFunc("/auth/spotify/login", oauthSpotifyLogin)
	mux.HandleFunc("/auth/spotify/callback", oauthSpotifyCallback)

	return mux
}

// TODO: Need to make sure we generate a single client and single auth url
func oauthSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// TODO: get data here
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
