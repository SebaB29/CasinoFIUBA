#!/bin/bash

set +e  # <-- NO cortar el script si algo falla
set -o pipefail  # <-- Pero sí reportar errores en pipes

BASE_URL="http://localhost:8080"
APUESTA=1000
TOTAL_PARTIDAS=5

echo "🧪 Iniciando test de $TOTAL_PARTIDAS partidas del juego de los vasos..."

# 🔐 Registro del usuario (ignora errores si ya existe)
echo "🔐 Registrando usuario de prueba..."
curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Tester",
    "apellido": "Vasos",
    "fecha_nacimiento": "1990-01-01",
    "email": "vasos@example.com",
    "password": "123456"
  }' > /dev/null

# 🔑 Login
echo "🔑 Logueando usuario..."
TOKEN=$(curl -s -X POST $BASE_URL/usuarios/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "vasos@example.com",
    "password": "123456"
  }' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "❌ No se pudo obtener el token."
  exit 1
fi

echo "✅ Token obtenido."

# 📊 Estadísticas
GANADAS=0
PERDIDAS=0

for ((i=1; i<=TOTAL_PARTIDAS; i++)); do
  echo ""
  echo "====================== 🧪 PARTIDA $i ======================"

  # Crear partida
  RESPUESTA_CREAR=$(curl -s -X POST $BASE_URL/vasos/nueva \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"apuesta\": $APUESTA}")

  ID=$(echo "$RESPUESTA_CREAR" | jq -r '.id')
  if [ -z "$ID" ] || [ "$ID" == "null" ]; then
    echo "❌ Error creando partida: $RESPUESTA_CREAR"
    continue
  fi

  echo "🎮 Partida creada - ID: $ID"

  # Elegir vaso aleatoriamente
  ELECCION=$(( RANDOM % 3 ))
  echo "📬 Jugada con elección: $ELECCION"

  RESPUESTA_JUGAR=$(curl -s -X POST $BASE_URL/vasos/jugar \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"id_partida\": $ID, \"eleccion\": $ELECCION}")

  echo "$RESPUESTA_JUGAR" | jq .

  ESTADO=$(echo "$RESPUESTA_JUGAR" | jq -r '.estado')

  if [ "$ESTADO" == "ganada" ]; then
    ((GANADAS++))
  elif [ "$ESTADO" == "perdida" ]; then
    ((PERDIDAS++))
  else
    echo "⚠️ Partida $ID no retornó estado válido."
  fi
done

# 📈 Resultados
echo ""
echo "====================== 📊 RESULTADOS ======================"
echo "✅ Partidas ganadas: $GANADAS"
echo "❌ Partidas perdidas: $PERDIDAS"
EXITO=$(echo "scale=2; $GANADAS*100/$TOTAL_PARTIDAS" | bc)
echo "🎯 Porcentaje de éxito: $EXITO%"
echo "==========================================================="
