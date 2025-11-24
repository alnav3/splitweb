# CSS Build stage
FROM node:18-alpine AS css-builder

WORKDIR /app

RUN npm init -y
RUN npm install tailwindcss@^3.4.0 @catppuccin/tailwindcss

COPY tailwind.config.js ./
COPY front/setup.css ./front/
COPY templates/ ./templates/
COPY main.go ./

RUN echo '@import "tailwindcss/base"; @import "tailwindcss/components"; @import "tailwindcss/utilities";' > ./front/docker.css

RUN npx tailwindcss -i ./front/docker.css -o ./style/tailwind.css --minify

# Go build stage
FROM golang:1.24.6-alpine AS go-builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=css-builder /app/style/tailwind.css ./style/

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=go-builder /app/main .

COPY --from=go-builder /app/style ./style

COPY --from=go-builder /app/db ./db

EXPOSE 8080

ENV GIN_MODE=release

CMD ["./main"]
