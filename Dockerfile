# Build stage
FROM node:18-alpine AS node-builder

WORKDIR /app

RUN npm init -y
RUN npm install tailwindcss @tailwindcss/cli @catppuccin/tailwindcss

COPY tailwind.config.js ./
COPY front/setup.css ./front/

# Create a simplified CSS input file for Docker build that doesn't rely on relative imports
RUN echo "@import 'tailwindcss';" > ./front/docker.css

COPY templates/ ./templates/
COPY main.go ./

RUN PATH=$PATH:./node_modules/.bin npx tailwindcss -i ./front/docker.css -o ./style/tailwind.css --minify

# Go build stage
FROM golang:1.24.6-alpine AS go-builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=node-builder /app/style/tailwind.css ./style/

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=go-builder /app/main .

COPY --from=go-builder /app/style ./style

COPY --from=go-builder /app/db ./db

EXPOSE 8080

ENV GIN_MODE=release

# Run the binary
CMD ["./main"]
