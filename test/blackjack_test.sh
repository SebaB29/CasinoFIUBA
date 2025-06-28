#!/bin/bash

set +e
echo "🧪 Iniciando test realista de Blackjack..."
BASE_URL="http://localhost:8080"

# Registro de usuario de prueba
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
  echo "❌ No se pudo obtener el token."
  exit 1
fi
echo "✅ Token obtenido."

# Contadores
GANADAS=0; PERDIDAS=0; EMPATADAS=0; RENDIDAS=0
SEGUROS=0; SPLITS=0; DOBLADAS=0
ERRORES=0; EN_CURSO=0

ACCIONES=("hit" "stand" "doblar" "rendirse" "seguro" "split")

for i in $(seq 1 20); do
  echo ""
  echo "====================== 🧪 PARTIDA $i ======================"

  RESPUESTA_CREAR=$(curl -s -X POST $BASE_URL/blackjack/nueva \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"apuesta": 100}')

  ID=$(echo "$RESPUESTA_CREAR" | jq -r '.id // empty')

  if [ -z "$ID" ]; then
    echo "❌ Error creando partida:"
    echo "$RESPUESTA_CREAR" | jq .
    ((ERRORES++))
    continue
  fi
  echo "🎮 Partida creada - ID: $ID"

  ACCION=${ACCIONES[$RANDOM % ${#ACCIONES[@]}]}
  echo "🏳️ Acción: ${ACCION^^}"

  if [[ "$ACCION" == "hit" ]]; then
    VALOR=0
    while [ $VALOR -lt 17 ]; do
      RESPUESTA_HIT=$(curl -s -X POST $BASE_URL/blackjack/hit \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"id_partida\": $ID}")
      VALOR=$(echo "$RESPUESTA_HIT" | jq -r '.valor // 99')
      [[ "$RESPUESTA_HIT" == *"error"* ]] && break
    done
    curl -s -X POST $BASE_URL/blackjack/stand \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID}" > /dev/null

  elif [[ "$ACCION" == "seguro" ]]; then
    RESPUESTA_ACCION=$(curl -s -X POST $BASE_URL/blackjack/seguro \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID}")
    if [[ "$RESPUESTA_ACCION" == *"seguro_pagado"* ]]; then
      ((SEGUROS++))
    elif [[ "$RESPUESTA_ACCION" == *"no se puede usar seguro"* ]]; then
      echo "ℹ️ Seguro no válido (banca sin As). Acción ignorada."
    else
      echo "ℹ️ Seguro ejecutado, partida continúa."
    fi

  else
    RESPUESTA_ACCION=$(curl -s -X POST $BASE_URL/blackjack/${ACCION} \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_partida\": $ID}")
    [[ "$RESPUESTA_ACCION" == *"doblado"* ]] && ((DOBLADAS++))
  fi

  # Obtener estado final
  RESULTADO=$(curl -s -X GET $BASE_URL/blackjack/estado/$ID \
    -H "Authorization: Bearer $TOKEN")
  echo "📋 Resultado final:"
  echo "$RESULTADO" | jq .

  ESTADO=$(echo "$RESULTADO" | jq -r '.estado // empty')
  CARTAS_USER=$(echo "$RESULTADO" | jq -c '.cartas_mano_1 // []')
  CARTAS_BANCA=$(echo "$RESULTADO" | jq -c '.cartas_banca // []')
  CARTAS_M2=$(echo "$RESULTADO" | jq -c '.cartas_mano_2 // []')

  echo "🃏 User: $CARTAS_USER | 🎴 Banca: $CARTAS_BANCA | ⚔️ Resultado: $ESTADO"

  if [ "$CARTAS_M2" != "[]" ]; then
    ((SPLITS++))
  fi

  case "$ESTADO" in
    "ganada")   ((GANADAS++)) ;;
    "perdida")  ((PERDIDAS++)) ;;
    "empatada") ((EMPATADAS++)) ;;
    "rendida")  ((RENDIDAS++)) ;;
    "en_curso")
      if [[ "$ACCION" == "split" || "$ACCION" == "seguro" ]]; then
        echo "ℹ️ Partida en curso tras $ACCION (esperado)."
      else
        echo "❌ Error: Partida quedó en curso. Acción: $ACCION"
        ((EN_CURSO++))
      fi ;;
    *)
      echo "❌ Error desconocido. Estado: $ESTADO"
      ((ERRORES++)) ;;
  esac
done

# Resultados finales
echo ""
echo "====================== 📊 RESULTADOS ======================"
echo "✅ Ganadas: $GANADAS"
echo "❌ Perdidas: $PERDIDAS"
echo "➖ Empatadas: $EMPATADAS"
echo "🏁 Rendidas: $RENDIDAS"
echo "💵 Seguros usados: $SEGUROS"
echo "✂️ Splits usados: $SPLITS"
echo "💥 Dobladas: $DOBLADAS"
echo "⚠️ Partidas en curso (error): $EN_CURSO"
echo "🚫 Otros errores: $ERRORES"
echo "==========================================================="

if [ $EN_CURSO -gt 0 ] || [ $ERRORES -gt 0 ]; then
  echo "❌ Test fallido: hay partidas no finalizadas correctamente."
  exit 1
else
  echo "✅ Todas las partidas finalizaron correctamente y con estados válidos."
  exit 0
fi
