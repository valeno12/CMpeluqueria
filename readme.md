## 🛠 Instalación y configuración

### 🔹 Clonar el repositorio  
```sh
git clone https://github.com/valeno12/CMpeluqueria.git
cd CMpeluqueria
```

### 🔹 Crear el archivo `.env`  
El proyecto necesita un archivo `.env` en la raíz con las credenciales de la base de datos.  

#### ✏️ **Ejemplo de `.env`:**
```
DB_USER=root
DB_PASSWORD=2328
DB_SERVER=db
DB_PORT=3306
DB_NAME=peluqueria

MYSQL_ROOT_PASSWORD=2328
MYSQL_DATABASE=peluqueria
MYSQL_USER=root
MYSQL_PASSWORD=2328
MYSQL_HOST=db
```

### 🔹 Levantar el proyecto con Docker  
```sh
sudo docker compose up -d --build
```

Esto construirá las imágenes y levantará los contenedores en segundo plano.

---

## 🛠 Solución de errores

### ❌ Error de acceso a la base de datos (`Access denied for user`)
Si ves un error similar a:
```
Error 1045 (28000): Access denied for user 'root'@'...' (using password: YES)
```
Asegurate de que en el archivo `.env` ambos valores (`DB_USER` y `MYSQL_USER`) sean `root` y que las contraseñas coincidan.

Si el error persiste, probá bajar los contenedores con:
```sh
sudo docker compose down -v
```
y luego volver a levantarlos con:
```sh
sudo docker compose up -d --build
```

