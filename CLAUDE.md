# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple Go web application that implements a TODO list functionality. The application uses Go's standard library for HTTP handling and HTML templating.

## Architecture

The application follows a simple monolithic architecture:
- `main.go`: Contains all server logic including route handlers and the main server setup
- `templates/todo.html`: HTML template for the TODO list interface using Go's template syntax
- `static/`: Contains static assets (CSS files) served directly

Key components:
- HTTP server running on port 8080
- Template-based rendering using Go's `html/template` package
- In-memory storage of TODO items (no database)
- Static file serving for CSS assets

## Development Commands

```bash
# Run the application
go run main.go

# Build the application
go build main.go

# Format Go code
go fmt main.go
```

## Routes

- `/todo` - Display the TODO list
- `/add` - Add a new TODO item (POST)
- `/static/` - Serve static files from the static directory

## Notes

- The application stores TODO items in memory, so data is lost on restart
- Default TODO items are initialized in Japanese: "洗顔", "朝食", "歯磨き"
- No authentication or persistent storage implemented