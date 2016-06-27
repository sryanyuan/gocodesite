package gocodecc

import (
	"github.com/gorilla/sessions"
)

const cookieKey = "gocodecc-session-store"

var store = sessions.NewCookieStore(cookieKey)
