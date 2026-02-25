# Classic Car Parts Search

A fuzzy search engine for classic car parts, powered by Google Sheets.

## Features

- **Authentication**: Login system with credentials stored in Google Sheets
- **Fuzzy Search**: Find parts even with typos or partial names using Levenshtein distance
- **Filters**: Filter by brand (Marca) and type (Tipo)
- **Estado Logic**: Items with "eliminado" in Estado column are hidden from search
- **Real-time Results**: Instant search with relevance scoring
- **Petrolhead UI**: Classic garage aesthetic with warm colors and retro typography
- **Mock Mode**: Test without connecting to Google Sheets

## Prerequisites

- Go 1.21 or higher
- Google Cloud account with Sheets API enabled (not required for mock mode)
- Service Account credentials (not required for mock mode)

## Quick Start (Mock Mode)

Run the application without Google Sheets using mock data:

```bash
# Install dependencies
go mod tidy

# Run in mock mode
MOCK_MODE=true go run cmd/server/main.go
```

Default mock credentials:
- Username: `admin`
- Password: `admin123`

The application will be available at `http://localhost:8080`

## Setup (Production Mode)

### 1. Google Cloud Configuration

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the **Google Sheets API**:
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Sheets API"
   - Click "Enable"

4. Create a Service Account:
   - Go to "IAM & Admin" > "Service Accounts"
   - Click "Create Service Account"
   - Assign a name and description
   - Grant "Editor" role

5. Generate Keys:
   - Click on the created service account
   - Go to "Keys" tab
   - Click "Add Key" > "Create new key"
   - Select JSON format
   - Download as `service-account.json`

### 2. Configure Google Sheet

Create a Google Sheet with two tabs:

#### Tab 1: "Repuestos" (Parts)
Columns:
- ID
- Nombre
- Marca
- Tipo
- Modelo
- Año
- Precio
- Descripcion
- Estado (leave empty to show, or "eliminado" to hide)

#### Tab 2: "Usuarios" (Users)
Columns:
- Usuario
- Password

Example:
| Usuario | Password |
|---------|----------|
| admin   | secret123|
| user    | password |

### 3. Share the Sheet

1. Click "Share"
2. Add the Service Account email (found in `service-account.json`)
3. Grant "Editor" permissions

### 4. Get Spreadsheet ID

Copy the Spreadsheet ID from the URL:
```
https://docs.google.com/spreadsheets/d/SPREADSHEET_ID_HERE/edit
```

### 5. Project Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your configuration:
   ```
   SPREADSHEET_ID=your-spreadsheet-id
   CREDENTIALS_PATH=service-account.json
   PORT=8080
   MOCK_MODE=false
   ```

3. Place your `service-account.json` in the project root

### 6. Run the Application

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`

## Usage

1. Open the application in your browser
2. Login with credentials from the "Usuarios" sheet (or mock credentials)
3. Enter a part name in the search box (fuzzy matching is enabled)
4. Use filters to narrow down results:
   - Brand (Marca)
   - Type (Tipo)

## Estado Column Behavior

- **Empty or no value**: Part is visible in search results
- **"eliminado"** (case-insensitive): Part is hidden from search results

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MOCK_MODE` | Use mock data instead of Google Sheets | `false` |
| `SPREADSHEET_ID` | Google Sheets ID | Required if not mock mode |
| `CREDENTIALS_PATH` | Path to service account JSON | `service-account.json` |
| `PORT` | Server port | `8080` |

## Project Structure

```
classicCarSearch/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── auth.go              # Authentication handlers
│   │   ├── handlers.go          # HTTP handlers
│   │   └── templates.go         # Template rendering
│   ├── models/
│   │   └── part.go              # Data models
│   └── services/
│       ├── auth.go              # Authentication service
│       ├── mock.go              # Mock data provider
│       ├── provider.go          # Data provider interface
│       ├── search.go            # Fuzzy search service
│       ├── session.go           # Session management
│       └── sheets.go            # Google Sheets service
├── static/
│   └── css/
│       └── style.css            # Petrolhead-style CSS
├── templates/
│   ├── index.html               # Main template
│   └── login.html               # Login template
├── .env.example                 # Environment template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Design

The UI follows a **petrolhead/garage aesthetic**:
- **Colors**: Cream, British Green, Rust Red, Steel Gray
- **Typography**: Playfair Display (headings), Merriweather (body)
- **Cards**: Beige/cream backgrounds with rounded corners
- **Buttons**: Subtle relief (light skeuomorphism)
- **Filters**: Dial-style selectors reminiscent of classic car radios

## Dependencies

- [quire](https://github.com/elbader17/quire) - Google Sheets database interface
- [fuzzysearch](https://github.com/lithammer/fuzzysearch) - Fuzzy string matching

## License

MIT
