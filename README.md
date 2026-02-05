# PDF Renderer Previewer

A Go template previewer tool for rendering and previewing HTML templates with dynamic data. Useful for designing PDF-ready HTML documents.

## Project Structure

```
pdf-renderer-previewer/
├── frontend/          # Vite-powered web UI
│   └── index.html     # Template editor & preview
├── render-api/        # Go backend API
│   ├── main.go        # HTTP server for template rendering
│   └── go.mod
└── templates/         # Sample HTML templates
    └── data-application.html
```

## Prerequisites

- **Go** 1.21+
- **Node.js** 20.19+ or 22.12+
- **npm**

## Getting Started

### 1. Start the Backend (Go API)

```bash
cd render-api
go mod init render-api   # Only needed first time
go run .
```

The API will start on `http://localhost:8080`.

### 2. Start the Frontend (Vite)

```bash
cd frontend
npm install
npx vite
```

The UI will be available at `http://localhost:5173`.

## Usage

1. Open `http://localhost:5173` in your browser
2. Enter your Go HTML template in the top-left textarea
3. Enter JSON data in the bottom-left textarea
4. The preview will auto-render on the right side

### Template Syntax

Uses Go's `html/template` syntax:

```html
<h1>Contract {{.contract_number}}</h1>
<p>Name: {{.customer.name}}</p>
```

With JSON data:

```json
{
  "contract_number": "FIN-2024-001",
  "customer": { "name": "John Doe" }
}
```

## API Reference

### POST /render/html

Renders a Go template with provided data.

**Request:**

```json
{
  "template": "<h1>Hello {{.name}}</h1>",
  "data": { "name": "World" }
}
```

**Response:**

```json
{
  "html": "<h1>Hello World</h1>"
}
```

**Error Response:**

```json
{
  "error": "template parse error: ..."
}
```

## Uploading Templates

To upload an HTML template to the document service:

```bash
curl --location 'https://microservices.sit.bravo.bfi.co.id/document/v1/document' \
--header 'api-secret: <API_SECRET>' \
--form 'file=@"templates/data-application.html"' \
--form 'ref_id="<UUID>"' \
--form 'id_type="test"' \
--form 'document_type="html"' \
--form 'source_system="LORA"' \
--form 'document_sequence="1"'
```

## Development

### CORS Configuration

The backend allows CORS from `localhost:5173` and `127.0.0.1:5173` for local development. Update [render-api/main.go](render-api/main.go#L80-L81) if using a different port.

## License

Internal use only.
