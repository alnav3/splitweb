module git.alnav.dev/alnav3/splitweb

go 1.22.4

require front v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

require (
	git.alnav.dev/alnav3/splitweb/api/db v0.0.0
	git.alnav.dev/alnav3/splitweb/api/encryption v0.0.0
	github.com/a-h/templ v0.2.731
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/gorilla/sessions v1.3.0
	github.com/libsql/sqlite-antlr4-parser v0.0.0-20240327125255-dbf53b6cbf06 // indirect
	github.com/tursodatabase/libsql-client-go v0.0.0-20240628122535-1c47b26184e8 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	nhooyr.io/websocket v1.8.10 // indirect
)

replace git.alnav.dev/alnav3/splitweb/api/db => ./api/db

replace git.alnav.dev/alnav3/splitweb/api/encryption => ./api/encryption

replace front => ./front
