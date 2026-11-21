# Kiln & Co.

A small, static ecommerce website for a handmade ceramics shop, written in Go using only the standard library (`net/http` and `html/template`). There are no external dependencies and no database.

## Features

- Landing page with featured products
- Shop catalogue and individual product pages
- Story page
- Enquiries form with a confirmation message
- Responsive layout in a single CSS file

## Routes

| Route          | Page                    |
|----------------|-------------------------|
| `/`            | Landing page            |
| `/shop`        | Product catalogue       |
| `/shop/{id}`   | Product detail          |
| `/story`       | About the brand         |
| `/enquiries`   | Contact / enquiry form  |

## Requirements

- Go 1.22 or newer

## Running the server

From the project folder, run:

```bash
go run .
```

The server starts on port 8080. Open `http://localhost:8080` in your browser.

## Project structure

```
kilnco/
├── main.go          # routes, handlers and product data
├── go.mod
├── templates/       # base layout and page templates
└── static/
    └── style.css
```

## Adding a product

Add a new entry to the `products` slice in `main.go`:

```go
{"glaze-cup", "Glaze Cup", "Mugs", "Short description here.", "#b5835a", 18},
```

The fields are ID, name, category, description, colour and price. It appears in the shop automatically.