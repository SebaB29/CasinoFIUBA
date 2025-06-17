#!/bin/bash

set -e
echo "🧪 Iniciando test largo de flujo Buscaminas con minas de 1 a 24..."

BASE_URL="http://localhost:8080"

mostrar_tablero() {
  local TABLERO_JSON=$1
  echo ""
  echo "🎲 Estado del tablero:"
  for Y in {0..4}; do
    for X in {0..4}; do
      CELDA=$(echo "$TABLERO_JSON" | jq -r ".[] | select(.x==$X and .y==$Y)")
      ABIERTA=$(echo "$CELDA" | jq -r ".abierta")
      if [ "$ABIERTA" == "true" ]; then
        printf "🟩 "
      else
        printf "⬛ "
      fi
    done
    echo ""
  done
  echo ""
}

# Registro
echo "🔐 Registrando usuario..."
curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Tester",
    "apellido": "Script",
    "fecha_nacimiento": "1990-01-01",
    "email": "tester@example.com",
    "password": "123456"
  }' > /dev/null || true

# Login
echo "🔑 Logueando usuario..."
TOKEN=$(curl -s -X POST $BASE_URL/usuarios/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tester@example.com",
    "password": "123456"
  }' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "❌ Error: No se pudo obtener el token."
  exit 1
fi
echo "✅ Token obtenido."

# Bucle de partidas
for MINAS in {1..24}; do
  echo ""
  echo "=========================== 🧱 PARTIDA #$MINAS - Minas: $MINAS ==========================="
  
  RESPUESTA_CREAR=$(curl -s -X POST $BASE_URL/buscaminas/nueva \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"minas\": $MINAS, \"apuesta\": 1000}")

  ID_PARTIDA=$(echo "$RESPUESTA_CREAR" | jq -r '.id_partida')
  echo "🆔 Partida ID: $ID_PARTIDA"
  TABLERO=$(echo "$RESPUESTA_CREAR" | jq '.tablero')
  mostrar_tablero "$TABLERO"

  VIVO=true
  CANT_ABIERTAS=0

  for X in {0..4}; do
    for Y in {0..4}; do
      if [ "$VIVO" = false ] || [ "$CANT_ABIERTAS" -ge 5 ]; then
        break
      fi

      echo "🧨 Abriendo celda ($X,$Y)..."
      RESPUESTA_ABRIR=$(curl -s -X POST $BASE_URL/buscaminas/abrir \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"id_partida\": $ID_PARTIDA, \"x\": $X, \"y\": $Y}")
      
      echo "📬 Respuesta:"
      echo "$RESPUESTA_ABRIR" | jq .

      # Manejo si se pisa una mina
      if echo "$RESPUESTA_ABRIR" | jq -e 'has("error")' > /dev/null; then
        echo "💥 ¡Boom! Partida perdida: $(echo "$RESPUESTA_ABRIR" | jq -r '.error')"
        VIVO=false
        break
      fi

      ESTADO=$(echo "$RESPUESTA_ABRIR" | jq -r '.estado')
      TABLERO=$(echo "$RESPUESTA_ABRIR" | jq '.tablero')
      mostrar_tablero "$TABLERO"

      if [ "$ESTADO" == "perdida" ]; then
        echo "💥 Partida perdida. No se puede seguir abriendo."
        VIVO=false
      else
        CANT_ABIERTAS=$((CANT_ABIERTAS + 1))
      fi
    done
    [ "$VIVO" = false ] && break
  done

  if [ "$VIVO" = true ]; then
    echo "🏃‍♂️ Retirándose..."
    RESPUESTA_RETIRO=$(curl -s -X POST $BASE_URL/buscaminas/retirarse \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID_PARTIDA}")
    
    echo "💰 Resultado del retiro:"
    echo "$RESPUESTA_RETIRO" | jq .
  else
    echo "🚫 No se retira porque perdió."
  fi

  echo "🕵️‍♂️ Mostrando minas en modo debug..."
  DEBUG_MINAS=$(curl -s -X GET $BASE_URL/buscaminas/debug/$ID_PARTIDA \
    -H "Authorization: Bearer $TOKEN")

  echo "📦 Respuesta cruda del debug:"
  echo "$DEBUG_MINAS" | jq .

  echo "💣 Ubicación de minas:"
  for Y in {0..4}; do
    for X in {0..4}; do
      if echo "$DEBUG_MINAS" | jq -e ".minas[] | select(.x==$X and .y==$Y)" > /dev/null; then
        printf "💣 "
      else
        printf "⬛ "
      fi
    done
    echo ""
  done
done

echo ""
echo "✅ Test largo finalizado correctamente."
#!/bin/bash

set -e
echo "🧪 Iniciando test largo de flujo Buscaminas con minas de 1 a 24..."

BASE_URL="http://localhost:8080"

mostrar_tablero() {
  local TABLERO_JSON=$1
  echo ""
  echo "🎲 Estado del tablero:"
  for Y in {0..4}; do
    for X in {0..4}; do
      CELDA=$(echo "$TABLERO_JSON" | jq -r ".[] | select(.x==$X and .y==$Y)")
      ABIERTA=$(echo "$CELDA" | jq -r ".abierta")
      if [ "$ABIERTA" == "true" ]; then
        printf "🟩 "
      else
        printf "⬛ "
      fi
    done
    echo ""
  done
  echo ""
}

# Registro
echo "🔐 Registrando usuario..."
curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Tester",
    "apellido": "Script",
    "fecha_nacimiento": "1990-01-01",
    "email": "tester@example.com",
    "password": "123456"
  }' > /dev/null || true

# Login
echo "🔑 Logueando usuario..."
TOKEN=$(curl -s -X POST $BASE_URL/usuarios/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tester@example.com",
    "password": "123456"
  }' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "❌ Error: No se pudo obtener el token."
  exit 1
fi
echo "✅ Token obtenido."

# Bucle de partidas
for MINAS in {1..24}; do
  echo ""
  echo "=========================== 🧱 PARTIDA #$MINAS - Minas: $MINAS ==========================="
  
  RESPUESTA_CREAR=$(curl -s -X POST $BASE_URL/buscaminas/nueva \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"minas\": $MINAS, \"apuesta\": 1000}")

  ID_PARTIDA=$(echo "$RESPUESTA_CREAR" | jq -r '.id_partida')
  echo "🆔 Partida ID: $ID_PARTIDA"
  TABLERO=$(echo "$RESPUESTA_CREAR" | jq '.tablero')
  mostrar_tablero "$TABLERO"

  VIVO=true
  CANT_ABIERTAS=0

  for X in {0..4}; do
    for Y in {0..4}; do
      if [ "$VIVO" = false ] || [ "$CANT_ABIERTAS" -ge 5 ]; then
        break
      fi

      echo "🧨 Abriendo celda ($X,$Y)..."
      RESPUESTA_ABRIR=$(curl -s -X POST $BASE_URL/buscaminas/abrir \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"id_partida\": $ID_PARTIDA, \"x\": $X, \"y\": $Y}")
      
      echo "📬 Respuesta:"
      echo "$RESPUESTA_ABRIR" | jq .

      # Manejo si se pisa una mina
      if echo "$RESPUESTA_ABRIR" | jq -e 'has("error")' > /dev/null; then
        echo "💥 ¡Boom! Partida perdida: $(echo "$RESPUESTA_ABRIR" | jq -r '.error')"
        VIVO=false
        break
      fi

      ESTADO=$(echo "$RESPUESTA_ABRIR" | jq -r '.estado')
      TABLERO=$(echo "$RESPUESTA_ABRIR" | jq '.tablero')
      mostrar_tablero "$TABLERO"

      if [ "$ESTADO" == "perdida" ]; then
        echo "💥 Partida perdida. No se puede seguir abriendo."
        VIVO=false
      else
        CANT_ABIERTAS=$((CANT_ABIERTAS + 1))
      fi
    done
    [ "$VIVO" = false ] && break
  done

  if [ "$VIVO" = true ]; then
    echo "🏃‍♂️ Retirándose..."
    RESPUESTA_RETIRO=$(curl -s -X POST $BASE_URL/buscaminas/retirarse \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID_PARTIDA}")
    
    echo "💰 Resultado del retiro:"
    echo "$RESPUESTA_RETIRO" | jq .
  else
    echo "🚫 No se retira porque perdió."
  fi

  echo "🕵️‍♂️ Mostrando minas en modo debug..."
  DEBUG_MINAS=$(curl -s -X GET $BASE_URL/buscaminas/debug/$ID_PARTIDA \
    -H "Authorization: Bearer $TOKEN")

  echo "📦 Respuesta cruda del debug:"
  echo "$DEBUG_MINAS" | jq .

  echo "💣 Ubicación de minas:"
  for Y in {0..4}; do
    for X in {0..4}; do
      if echo "$DEBUG_MINAS" | jq -e ".minas[] | select(.x==$X and .y==$Y)" > /dev/null; then
        printf "💣 "
      else
        printf "⬛ "
      fi
    done
    echo ""
  done
done

echo ""
echo "✅ Test largo finalizado correctamente."
