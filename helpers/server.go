package helpers

// Go server to handle REST operations for Spotify API

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

func server() {
  server := &http.Server{
    Addr: fmt.Sprintf(":8080"),
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

func oauthSpotifyLogin(w http.ResponseWriter, r *http.Request) {
  oauthState := generateStateOAuthCookie(w) 
}

func generateStateOAuthCookie(w http.ResponseWriter) string {
  var expiration = time.Now().Add(365 * 24 * time.Hour)

  bitStr := make([]byte, 16)
  rand.Read(bitStr)
  state := base64.URLEncoding.EncodeToString(bitStr)
  cookie := http.Cookie{name: "spotifyoauthstate", Value: state, Expires: expiration}
  http.SetCookie(w, &cookie)

  return state
}
