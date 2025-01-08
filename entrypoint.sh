#!/bin/sh

# Esperar a que la base de datos esté disponible
echo "Esperando a que la base de datos esté disponible..."
until nc -z -v -w30 db 3306
do
  echo "Esperando por DB..."
  sleep 5
done

# Ejecutar migraciones y seeders si es necesario
echo "Ejecutando migraciones..."
./main setup

# Iniciar la aplicación
echo "Iniciando la aplicación..."
exec "$@"
