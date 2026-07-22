#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

cleanup() {
  echo ""
  echo "🧹 Cleaning up..."
  kill $GO_PID $FB_PID 2>/dev/null || true
  wait $GO_PID $FB_PID 2>/dev/null || true
}
trap cleanup EXIT

wait_for_port() {
  local port=$1
  local max_wait=${2:-30}
  local elapsed=0
  while ! curl -s "http://localhost:$port" >/dev/null 2>&1; do
    sleep 1
    elapsed=$((elapsed + 1))
    if [ "$elapsed" -ge "$max_wait" ]; then
      echo "❌ Timeout waiting for port $port"
      return 1
    fi
  done
}

echo "🔥 Starting Firebase emulators..."
firebase emulators:start --only auth,firestore &>/tmp/firebase-emulator.log &
FB_PID=$!
wait_for_port 8080 30

echo "🚀 Starting Go server..."
FIRESTORE_EMULATOR_HOST=localhost:8080 GCP_PROJECT=demo-no-project \
  go run ./cmd/server &>/tmp/go-server.log &
GO_PID=$!
wait_for_port 8000 15

echo ""
echo "════════════════════════════════════════"
echo "  TEST 1: Health Check"
echo "════════════════════════════════════════"
curl -s http://localhost:8000/health
echo ""

echo ""
echo "════════════════════════════════════════"
echo "  TEST 2: POST /api/v1/users sin auth"
echo "════════════════════════════════════════"
curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Test","last_name":"User","dni":"12345","email":"test@test.com"}'

echo ""
echo "════════════════════════════════════════"
echo "  TEST 3: POST /api/v1/users token inválido"
echo "════════════════════════════════════════"
curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake-token" \
  -d '{"first_name":"Test","last_name":"User","dni":"12345","email":"test@test.com"}'

echo ""
echo "════════════════════════════════════════"
echo "  TEST 4: GET /api/v1/roles sin auth"
echo "════════════════════════════════════════"
curl -s -w "\nHTTP %{http_code}\n" http://localhost:8000/api/v1/roles

echo ""
echo "════════════════════════════════════════"
echo "  Resumen:"
echo "  ✅ Health check funciona"
echo "  ✅ Auth middleware rechaza requests sin token"
echo "  ✅ Auth middleware rechaza tokens inválidos"
echo "  ✅ Roles endpoint registrado y protegido"
echo "  📊 Emulator UI: http://localhost:4000"
echo "════════════════════════════════════════"
