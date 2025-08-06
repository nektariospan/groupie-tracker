# Groupie Tracker

🎸 **Groupie Tracker** is a Go-based web application that displays information about musical artists, including their members, first album, locations, concert dates, and more. It features a responsive, 3D flip-card UI and server-side pagination.

## Features

- **Go Backend**: Fetches data asynchronously from the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api).
- **HTML Templates**: Uses Go's `html/template` for clean server-side rendering.
- **Responsive UI**: CSS Grid layout adapts from 5 columns down to 1 on small screens.
- **3D Flip Cards**: Hover or click artists and team member cards to see details on the back.
- **Pagination**: `?limit=10|20|all` and `?page=N` query parameters for easy navigation.
- **Middleware**: Logging and recovery with standard `net/http` middleware.
- **Team Page**: Profiles of the development team with social links (Discord, GitHub).

## Getting Started

### Prerequisites

- Go 1.24+
- Internet connection (to fetch data from the API)

### Installation

1. Run the server:
   ```bash
   go run .
   ```

2. Open your browser at `http://localhost:8080`.

## Project Structure

```
GROUPIE-TRACKER/
├── data/
│   ├── fetch.go       # Data fetching logic
│   └── models/
│       └── types.go   # Domain models
├── static/
│   ├── css/           # Stylesheets
│   │   ├── artist.css
│   │   ├── base.css
│   │   ├── error.css
│   │   ├── style.css
│   │   └── team.css
│   └── images/        # Static images
├── templates/         # HTML templates
│   ├── index.html
│   ├── team.html
│   ├── artist.html
│   └── errors/
│       └── error.html
├── web/               # Web layer
│   ├── helpers.go     # Template rendering helpers
│   ├── middleware.go  # HTTP middleware (recovery)
│   └── server.go      # Routing and handlers
├── go.mod             # Go module definition
├── main.go            # Application entry point
└── README.md          # Project documentation
```

## API Endpoints

- `GET /artists` – List all artist summaries
- `GET /artist/{id}` – Artist details (members, albums, etc.)
- `GET /team` – Development team page
- `GET /static/...` – Static files (CSS, images)

## License

This project is licensed under the [MIT License](LICENSE).