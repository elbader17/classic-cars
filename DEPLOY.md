# Despliegue en Google Cloud Run + Google Sheets

## Prerrequisitos

- Cuenta de Google Cloud (gratis)
- Go instalado localmente
- Google Cloud CLI (`gcloud`) instalado

---

## Paso 1: Crear proyecto en Google Cloud

1. Ve a [Google Cloud Console](https://console.cloud.google.com/)
2. Crea un nuevo proyecto: `classic-car-search`
3. Habilita las APIs necesarias:

```bash
gcloud auth login
gcloud projects create classic-car-search --name="Classic Car Search"
gcloud config set project classic-car-search
```

Habilitar APIs:
```bash
gcloud services enable run.googleapis.com containerregistry.googleapis.com sheets.googleapis.com drive.googleapis.com
```

---

## Paso 2: Configurar Google Sheets

### 2.1 Crear la hoja de cálculo

1. Crea una hoja de cálculo en Google Sheets
2. Nombra la pestaña como: `Repuestos`
3. Define las columnas (fila 1):

| A | B | C | D | E | F | G | H |
|---|---|---|---|---|---|---|---|
| ID | Nombre | Marca | Tipo | Modelo | Año | Precio | Descripcion |
| 1 | Carburador Ford V8 | Ford | Motor | Mustang | 1967 | 450 | Carburador original |

**Columnas adicionales opcionales:**
- `Imagenes` - URLs separadas por coma
- `Estado` - "eliminado" para ocultar

### 2.2 Crear Service Account

```bash
gcloud iam service-accounts create classic-car-sa \
    --display-name="Classic Car Search Service Account"
```

### 2.3 Generar clave JSON

```bash
gcloud iam service-accounts keys create service-account.json \
    --iam-account=classic-car-sa@classic-car-search.iam.gserviceaccount.com
```

### 2.4 Compartir Sheets con Service Account

1. Abre tu Google Sheet
2. Comparte con el email del service account:
   ```
   classic-car-sa@classic-car-search.iam.gserviceaccount.com
   ```
3. Dale rol de **Editor**

### 2.5 Obtener el Spreadsheet ID

De la URL del Sheet:
```
https://docs.google.com/spreadsheets/d/【SPREADSHEET_ID】/edit
```

Copia el ID (texto entre `/d/` y `/edit`)

---

## Paso 3: Configurar el código

### 3.1 Actualizar go.mod con build info

```go
// go.mod
module github.com/eduardo/classicCarSearch
```

### 3.2 Modificar main.go para producción

Edita `cmd/server/main.go`:

```go
func main() {
    mockMode := getEnvBool("MOCK_MODE", false)  // Cambiar a false para producción
    
    // ...
}
```

### 3.3 Crear Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .
COPY service-account.json .

EXPOSE 8080

CMD ["./server"]
```

---

## Paso 4: Desplegar en Cloud Run

### 4.1 Build y deploy

```bash
# Configura region
gcloud config set run/region us-central1

# Build y deploy
gcloud run deploy classic-car-search \
    --source . \
    --service-account classic-car-sa@classic-car-search.iam.gserviceaccount.com \
    --allow-unauthenticated \
    --set-env-vars SPREADSHEET_ID=TU_SPREADSHEET_ID,MOCK_MODE=false
```

### 4.2 Obtener URL del servicio

Al finalizar verás algo como:
```
Service URL: https://classic-car-search-xxx.run.app
```

---

## Paso 5: Conectar frontend

### 5.1 Actualizar index.html

Edita la línea del `API_BASE` en `index.html`:

```javascript
const API_BASE = 'https://classic-car-search-xxx.run.app';
```

### 5.2 Hospedar frontend (opciones)

**Opción A: Cloud Run estático**
- Pon el `index.html` en una carpeta `public/`
- Actualiza el Dockerfile para servir archivos estáticos

**Opción B: Firebase Hosting**
```bash
firebase init hosting
firebase deploy
```

**Opción C: Google Cloud Storage**
```bash
gsutil mb -l us-central1 gs://classic-car-web
gsutil web set -m index.html gs://classic-car-web
gsutil cp index.html gs://classic-car-web/
```

---

## Variables de entorno en Cloud Run

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `SPREADSHEET_ID` | ID del Google Sheet | `1abc123...xyz` |
| `MOCK_MODE` | Usar datos mock | `false` |
| `PORT` | Puerto del servidor | `8080` |
| `CREDENTIALS_PATH` | Path del JSON | `service-account.json` |

---

## Solución de problemas

### Error: "Unable to verify credentials"
- Verifica que el service account tenga acceso al Sheet
- Confirma que compartiste el Sheet con el email del SA

### Error: "Sheet not found"
- Verifica que el `SPREADSHEET_ID` sea correcto
- Confirma que la pestaña se llame `Repuestos`

### Ver logs
```bash
gcloud logs read --service=classic-car-search --limit=50
```

---

## RESUMEN - URLs y Configuración

### Backend (Cloud Run)
- **URL:** `https://classic-car-search-983986360116.us-central1.run.app`
- **Proyecto:** classic-car-search
- **Service Account:** classic-car-sa@classic-car-search.iam.gserviceaccount.com
- **Spreadsheet ID:** 1x2C-AIm5iyYzy9duC2iYD5k_bHcqpjmeo9h4UmnrFeA

### Google Sheets
- **Hoja "Repuestos":** Datos de repuestos
- **Hoja "Usuarios":** Credenciales de acceso (opcional, hay fallback)

### Credenciales de acceso (fallback)
- **Usuario:** admin
- **Contraseña:** admin123

### Para actualizar el backend
```bash
# 1. Hacer cambios en el código
# 2. Compilar y subir
gcloud builds submit --tag gcr.io/classic-car-search/classic-car-search:latest .

# 3. Desplegar
gcloud run deploy classic-car-search \
    --image gcr.io/classic-car-search/classic-car-search:latest \
    --service-account classic-car-sa@classic-car-search.iam.gserviceaccount.com \
    --allow-unauthenticated \
    --set-env-vars SPREADSHEET_ID=1x2C-AIm5iyYzy9duC2iYD5k_bHcqpjmeo9h4UmnrFeA,MOCK_MODE=false \
    --update-secrets GOOGLE_CREDENTIALS=service-account-json-key:latest \
    --region us-central1
```

### Para hostear el frontend
1. Firebase Hosting:
   ```bash
   firebase init hosting
   firebase deploy
   ```
2. O cualquier hosting estático (Netlify, Vercel, etc.)
