# CONFIGURACIÓN DE CLASSIC CAR PARTS CON GOOGLE SHEETS

## 🚀 CONFIGURACIÓN RÁPIDA

### 1. CREAR GOOGLE SHEETS

1. **Crear una nueva hoja de cálculo en Google Sheets**
   - Ve a [sheets.google.com](https://sheets.google.com)
   - Crea una nueva hoja en blanco
   - ¡Nómina esta hoja "Repuestos"!

2. **Crear las columnas necesarias**
   En la hoja "Repuestos", crea las siguientes columnas en la fila 1:

   | ID | Nombre | Marca | Tipo | Modelo | Año | Precio | Descripcion | Estado | Imagenes |

   **Ejemplo de datos:**
   ```
   ID    | Nombre              | Marca    | Tipo        | Modelo  | Año | Precio | Descripcion              | Estado | Imagenes
   ----- | ------------------- | -------- | ----------- | ------- | ---- | ------ | ------------------------- | ------ | -------
   1     | Carburador Ford V8  | Ford     | Motor       | Mustang | 1967 | 450.00 | Carburador original       |        | https://imagen1.jpg,https://imagen2.jpg
   2     | Radiador Chevrolet  | Chevrolet| Enfriamiento| Camaro  | 1969 | 320.00 | Radiador reconstruido     |        | https://imagen.jpg
   3     | Bomba de agua Ford  | Ford     | Enfriamiento| Falcon  | 1970 | 150.00 | Nueva, en caja            |        | 
   ```

### 2. CREAR HOJA DE USUARIOS

3. **Crear una segunda hoja llamada "Usuarios"**
   - Haz clic en el botón "+" para agregar una nueva hoja
   - Nómbrala "Usuarios"
   - Crea estas columnas:

   | Usuario | Password |

   **Ejemplo:**
   ```
   Usuario | Password
   ------- | --------
   admin   | admin123
   usuario | test
   ```

### 3. CREAR SCRIPT DE GOOGLE APPS

4. **Abrir Google Apps Script**
   - En Google Sheets, ve a: `Extensiones → Apps Script`
   - Se abrirá una nueva pestaña con el editor de Apps Script

5. **Reemplazar el código por defecto**
   - Elimina todo el código que aparece por defecto
   - Copia y pega el contenido del archivo `google-sheets-api.gs`
   - Guarda el proyecto con `Ctrl+S` o `Cmd+S`
   - Nómbralo "Classic Car Parts API"

6. **Publicar el script como API**
   - Ve a: `Publicar → Implementar como API web`
   - En la ventana emergente, haz clic en "Implementar"
   - Selecciona "¡Nueva implementación!"
   - Configuración:
     - Descripción: "Classic Car Parts API"
     - Tipo de ejecución: `doGet`
     - Ámbito de acceso: `Cualquiera (incluso anónimo)`
   - Haz clic en "Implementar"
   - Acepta los permisos que solicita

7. **Copiar la URL de la API**
   - Después de implementar, copia la URL que aparece
   - Ejemplo: `https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec`
   - Reemplaza `YOUR_SCRIPT_ID` con tu ID real

### 4. CONFIGURAR EL HTML

8. **Editar el archivo HTML**
   - Abre el archivo `index.html`
   - Busca esta línea:
     ```javascript
     const API_BASE = 'https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec';
     ```
   - Reemplaza `YOUR_SCRIPT_ID` con tu ID real de Google Apps Script

9. **Probar la aplicación**
   - Abre el archivo `index.html` en tu navegador
   - Inicia sesión con las credenciales de la hoja "Usuarios"
   - ¡Listo! Deberías ver tus datos de Google Sheets

## 📋 ESTRUCTURA DE LAS HOJAS

### Hoja "Repuestos" (Obligatoria)
```
ID (string)    - Identificador único
Nombre (string) - Nombre del repuesto
Marca (string)  - Marca del fabricante
Tipo (string)   - Categoría (Motor, Enfriamiento, etc.)
Año (string)    - Año del modelo
Precio (number) - Precio en USD
Descripcion (string) - Descripción del repuesto
Estado (string) - "eliminado" para ocultar, vacío para mostrar
Imagenes (string) - URLs de imágenes separadas por comas
```

### Hoja "Usuarios" (Obligatoria)
```
Usuario (string) - Nombre de usuario
Password (string) - Contraseña
```

### Hoja "Sessions" (Automática)
- El script crea esta hoja automáticamente
- Almacena las sesiones activas
- No necesitas crearla manualmente

## 🔧 PERMISOS REQUERIDOS

El script solicitará estos permisos al implementarlo:
- ✅ Ver y administrar tus hojas de cálculo de Google Drive
- ✅ Conectar a API externas
- ✅ Publicar aplicaciones web

## 📊 CARACTERÍSTICAS IMPLEMENTADAS

### ✅ Autenticación
- Login con usuarios de Google Sheets
- Tokens de sesión con 2 horas de validez
- Logout y gestión de sesiones

### ✅ Búsqueda Fuzzy
- Búsqueda por nombre, marca, tipo y modelo
- Relevancia por coincidencias
- Tolerancia a errores de tipeo

### ✅ Filtros
- Filtrar por marca
- Filtrar por tipo/categoría
- Combinación de filtros y búsqueda

### ✅ Gestión de Estado
- Columna "Estado": "eliminado" oculta el repuesto
- Partes con estado vacío se muestran normalmente

### ✅ Galería de Imágenes
- Múltiples imágenes por repuesto
- Visor con zoom y miniaturas
- Imágenes por defecto si no hay URLs

## 🚨 DEPURACIÓN Y MANTENIMIENTO

### Probar la API manualmente
1. Abre la URL de tu API en el navegador
2. Añade: `?endpoint=parts&token=test`
3. Deberías ver los datos en formato JSON

### Limpiar sesiones
- En Google Sheets, ve a: `Extensiones → Classic Car Parts API → Admin → Clear Sessions`
- Esto elimina todas las sesiones activas

### Ver información de depuración
- Ve a: `Extensiones → Classic Car Parts API → Debug Info`
- Muestra conteos de partes, usuarios y sesiones

## 📋 TIPS ÚTILES

### Agregar nuevos repuestos
1. Añade filas en la hoja "Repuestos"
2. Asegúrate de completar todas las columnas
3. Deja "Estado" vacío para que se muestre
4. Separa URLs de imágenes con comas

### Actualizar usuarios
1. Edita la hoja "Usuarios"
2. Añade nuevos usuarios con contraseñas
3. Los cambios se reflejan inmediatamente

### Cambiar la apariencia
- El HTML ya incluye el diseño petrolhead
- Puedes modificar los estilos CSS en el archivo `index.html`
- Los colores y fuentes están definidos en la configuración de Tailwind

## 🎯 LISTO PARA USAR

¡Felicidades! Tu aplicación está lista para usar:

1. ✅ Google Sheets configurada
2. ✅ Script de Apps Script implementado
3. ✅ HTML modificado con tu API
4. ✅ Datos dinámicos desde Google Sheets

Abre `index.html` en tu navegador y disfruta de tu catálogo de repuestos clásicos conectado a Google Sheets.