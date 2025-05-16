#!/bin/bash

echo "🔄 Deteniendo contenedores anteriores (si existen)..."
docker-compose down

echo "🚀 Levantando el entorno de desarrollo..."
docker-compose up --build
