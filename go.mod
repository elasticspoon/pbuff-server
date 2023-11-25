module github.com/elasticspoon/pbuff-server

go 1.21.3

require (
	github.com/elasticspoon/pbuff-server/server v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.0.10
)

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace github.com/elasticspoon/pbuff-server/server => ./server
