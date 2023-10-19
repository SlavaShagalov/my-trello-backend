package middleware

import "net/http"

type Middleware func(h http.HandlerFunc) http.HandlerFunc
