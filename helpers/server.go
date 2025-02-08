package helpers

// Go server to handle REST operations for Spotify API

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"mikehaus/spoofify/views"
	"net/http"
	"time"

  "github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func InitServer(auth *SpotifyAuth) {
	r := mux.NewRouter()

  // Spotify Routes
	r.HandleFunc("/auth/spotify/login", auth.HandleSpotifyLogin).Methods("GET")
	r.HandleFunc("/auth/spotify/callback", auth.SpotifyAuthCallback).Methods("GET")

  // templ generated component routes
	r.Handle("/auth/spotify/init", templ.Handler(views.Login(auth.config.RedirectURL)))

	// start on secondary thread to ensure it's not blocking. UI will run on main thread
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()
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
