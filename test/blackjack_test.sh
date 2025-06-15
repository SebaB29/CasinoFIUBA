#!/bin/bash

set -e
echo "🧪 Iniciando test de 10 partidas de Blackjack..."

BASE_URL="http://localhost:8080"

# Registro (ignorar si ya existe)
echo "🔐 Registrando usuario de prueba..."
curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Tester", "apellido": "Blackjack", "email": "blackjack@example.com", "password": "123456", "fecha_nacimiento": "1990-01-01"}' > /dev/null || true

# Login
echo "🔑 Logueando usuario..."
TOKEN=$(curl -s -X POST $BASE_URL/usuarios/login \
  -H "Content-Type: application/json" \
  -d '{"email": "blackjack@example.com", "password": "123456"}' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Error: No se pudo obtener el token."
  exit 1
fi
echo "✅ Token obtenido."

GANADAS=0
PERDIDAS=0
EMPATADAS=0
ERRORES=0

for i in $(seq 1 10); do
  echo ""
  echo "====================== 🧪 PARTIDA $i ======================"

  # Crear partida
  RESPUESTA_CREAR=$(curl -s -X POST $BASE_URL/blackjack/nueva \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"apuesta": 100}')

  ID=$(echo "$RESPUESTA_CREAR" | jq -r '.id')

  if [ -z "$ID" ] || [ "$ID" = "null" ]; then
    echo "❌ Error creando partida:"
    echo "$RESPUESTA_CREAR" | jq .
    ERRORES=$((ERRORES + 1))
    continue
  fi

  echo "🎮 Partida creada - ID: $ID"

  # Mostrar estado inicial
  echo "📋 Estado inicial:"
  curl -s -X GET $BASE_URL/blackjack/estado/$ID \
    -H "Authorization: Bearer $TOKEN" | jq .

  # Decidir acción al azar: hit o stand
  if [ $((RANDOM % 2)) -eq 0 ]; then
    echo "👉 Hit 1"
    RESPUESTA_HIT=$(curl -s -X POST $BASE_URL/blackjack/hit \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID, \"accion\": \"hit\"}")
    
    echo "📋 Estado tras hit:"
    echo "$RESPUESTA_HIT" | jq .

    ESTADO=$(echo "$RESPUESTA_HIT" | jq -r '.estado // empty')
    if [ "$ESTADO" = "perdida" ]; then
      echo "🛑 La partida ya está finalizada con estado: $ESTADO"
      PERDIDAS=$((PERDIDAS + 1))
      continue
    fi
  fi

  echo "✋ Stand"
  RESPUESTA_FINAL=$(curl -s -X POST $BASE_URL/blackjack/stand \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"id_partida\": $ID, \"accion\": \"stand\"}")

  echo "$RESPUESTA_FINAL" | jq .

  ESTADO=$(echo "$RESPUESTA_FINAL" | jq -r '.estado // empty')

  case $ESTADO in
    "ganada") GANADAS=$((GANADAS + 1)) ;;
    "perdida") PERDIDAS=$((PERDIDAS + 1)) ;;
    "empatada") EMPATADAS=$((EMPATADAS + 1)) ;;
    *)
      echo "⚠️ Estado inesperado o error"
      ERRORES=$((ERRORES + 1))
      ;;
  esac
done

echo ""
echo "====================== 📊 RESULTADOS ======================"
echo "✅ Ganadas: $GANADAS"
echo "❌ Perdidas: $PERDIDAS"
echo "➖ Empatadas: $EMPATADAS"
echo "⚠️ Errores: $ERRORES"
echo "==========================================================="
