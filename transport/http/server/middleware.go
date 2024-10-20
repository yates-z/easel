package server

type Middleware func(HandlerFunc) HandlerFunc
