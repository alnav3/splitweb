# SplitWeb

A simple web application built with Go, Templ, TailwindCSS, and HTMX.

## Project Structure

- `main.go`: Entry point and HTTP server configuration
- `templates/`: Templ template files
- `front/`: Front-end source files (CSS)
- `style/`: Compiled static assets (CSS)

## Running the Application

1. Start the development server:
   ```
   air
   ```
   This uses the configuration in `.air.toml` for hot reloading.

2. Generate Templ files (after making changes to `.templ` files):
   ```
   templ generate
   ```

3. Build TailwindCSS:
   ```
   npx tailwindcss -i ./front/setup.css -o ./style/tailwind.css
   ```

4. Access the application at http://localhost:8080

## Features

- Go backend with standard library HTTP server
- Templ for HTML templating
- TailwindCSS for styling
- HTMX for interactive UI without JavaScript
- Alpine.js for client-side interactivity