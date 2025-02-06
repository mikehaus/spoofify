package views

import "github.com/a-h/templ"

func Login(redirectUri string) templ.Component {
  return login(redirectUri)
} 
