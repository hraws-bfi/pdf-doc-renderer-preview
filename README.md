# PDF Renderer Previewer

A Go template previewer tool for rendering and previewing HTML templates with dynamic data. Useful for designing PDF-ready HTML documents.

## Project Structure

```
pdf-renderer-previewer/
├── .env.example       # Environment variables template
├── frontend/          # Vite-powered web UI
│   └── index.html     # Template editor & preview
├── render-api/        # Go backend API
│   ├── main.go        # Entry point, server setup
│   ├── config.go      # Configuration loading
│   ├── types.go       # Request/response types
│   ├── handlers.go    # HTTP handlers
│   ├── middleware.go  # CORS middleware
│   ├── utils.go       # Utility functions
│   └── go.mod
└── templates/         # Saved HTML templates
```

## Prerequisites

- **Go** 1.21+
- **Node.js** 20.19+ or 22.12+
- **npm**

## Getting Started

### 1. Configure Environment

```bash
cp .env.example .env
# Edit .env with your settings
```

### 2. Start the Backend (Go API)

```bash
cd render-api
go run .
```

The API will start on `http://localhost:8080`.

### 3. Start the Frontend (Vite)

```bash
cd frontend
npm install
npx vite
```

The UI will be available at `http://localhost:5173`.

## Usage

1. Open `http://localhost:5173` in your browser
2. Enter your Go HTML template in the top editor
3. Enter JSON data in the bottom editor
4. The preview will auto-render on the right side
5. Save templates with versioning using the "Save Template" button
6. Upload templates to DMS using the "Upload to DMS" button

### Template Syntax

Uses Go's `html/template` syntax:

```html
<h1>Contract {{.contract_number}}</h1>
<p>Name: {{.customer.name}}</p>

{{if .show_details}}
  <p>Details are visible</p>
{{end}}

{{range .items}}
  <li>{{.name}}: {{.price}}</li>
{{end}}
```

With JSON data:

```json
{
  "contract_number": "FIN-2024-001",
  "customer": { "name": "John Doe" },
  "show_details": true,
  "items": [
    { "name": "Item 1", "price": 100 },
    { "name": "Item 2", "price": 200 }
  ]
}
```

## Environment Variables

Create a `.env` file in the project root:

```bash
# Allowed CORS origins (comma-separated)
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:5174

# DMS (Document Management Service) Configuration
DMS_API_URL=https://microservices.sit.bravo.bfi.co.id/document/v1/document
DMS_API_SECRET=your-api-secret-here
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

### POST /render/pdf

Renders a Go template to PDF and returns the PDF file for download.

**Prerequisites:** Chrome or Chromium must be installed on the system.

**Request:**

```json
{
  "template": "<h1>Hello {{.name}}</h1>",
  "data": { "name": "World" },
  "filename": "document"
}
```

**Response:** Binary PDF file with `Content-Type: application/pdf`

### POST /templates/save

Saves a template with automatic versioning.

**Request:**

```json
{
  "name": "data-application",
  "template": "<html>...</html>",
  "data": "{\"key\": \"value\"}"
}
```

**Response:**

```json
{
  "filename": "data-application-v1.html",
  "version": 1
}
```

### GET /templates/list

Lists all saved templates.

**Response:**

```json
{
  "templates": [
    { "name": "data-application", "filename": "data-application-v2.html", "version": 2 },
    { "name": "data-application", "filename": "data-application-v1.html", "version": 1 }
  ]
}
```

### POST /templates/upload-dms

Uploads a saved template to the Document Management Service.

**Request:**

```json
{
  "filename": "data-application-v1.html",
  "ref_id": "550e8400-e29b-41d4-a716-446655440000",
  "id_type": "template",
  "document_type": "html",
  "source_system": "LORA",
  "document_sequence": "1"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Template data-application-v1.html uploaded successfully",
  "response": "..."
}
```

## Manual Upload via cURL

To upload an HTML template to the document service manually:

```bash
curl --location 'https://microservices.sit.bravo.bfi.co.id/document/v1/document' \
--header 'api-secret: <API_SECRET>' \
--form 'file=@"templates/data-application-v1.html"' \
--form 'ref_id="<UUID>"' \
--form 'id_type="template"' \
--form 'document_type="html"' \
--form 'source_system="LORA"' \
--form 'document_sequence="1"'
```

## License

Internal use only.
