#!/bin/bash

BASE_URL="http://localhost:8000/api"
EMAIL="mati@example.com"
PASSWORD="123456"

# --- LOGIN ---
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

echo "Login response:"
echo "$LOGIN_RESPONSE" | jq .

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

# izvlačimo user_id iz JWT payload-a (ispravno base64url decode)
PAYLOAD=$(echo "$TOKEN" | cut -d "." -f2)

LEN=$(( ${#PAYLOAD} % 4 ))
if [ $LEN -eq 2 ]; then
  PAYLOAD="${PAYLOAD}=="
elif [ $LEN -eq 3 ]; then
  PAYLOAD="${PAYLOAD}="
fi

USER_ID=$(echo "$PAYLOAD" | tr '_-' '/+' | base64 --decode | jq -r '.user_id')

echo "Token: $TOKEN"
echo "Author ID: $USER_ID"

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ] || [ -z "$USER_ID" ] || [ "$USER_ID" == "null" ]; then
    echo "❌ Login failed, check credentials or API response."
    exit 1
fi

# --- CREATE TOUR ---
echo "Creating tour..."
TOUR_RESPONSE=$(curl -s -X POST "$BASE_URL/tours" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
        "name": "Planinska tura",
        "description": "Lepa pešačka tura kroz planine",
        "difficulty": "medium",
        "tags": "planine,priroda"
      }')

echo "$TOUR_RESPONSE" | jq .
TOUR_ID=$(echo "$TOUR_RESPONSE" | jq -r '.id')

echo "Tour ID: $TOUR_ID"

# --- ADD KEYPOINTS ---
echo "Adding keypoints..."
for kp in \
  '{"name":"Početna tačka","description":"Početak ture","latitude":44.8176,"longitude":20.4569,"order":1}' \
  '{"name":"Vidikovac","description":"Pogled na planine","latitude":44.8200,"longitude":20.4600,"order":2}' \
  '{"name":"Planinski vrh","description":"Najviša tačka ture","latitude":44.823,"longitude":20.465,"order":3}'; do
    RESPONSE=$(curl -s -X POST "$BASE_URL/tours/$TOUR_ID/keypoints" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "$kp")
    echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
done

# --- ADD DURATIONS ---
echo "Adding durations..."
for duration in \
  '{"transport":"walk","minutes":120}' \
  '{"transport":"bike","minutes":45}'; do
    RESPONSE=$(curl -s -X POST "$BASE_URL/tours/$TOUR_ID/durations" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "$duration")
    echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
done

# --- PUBLISH TOUR ---
echo "Publishing tour..."
RESPONSE=$(curl -s -X POST "$BASE_URL/tours/$TOUR_ID/publish" \
  -H "Authorization: Bearer $TOKEN")
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"

# --- ARCHIVE TOUR ---
echo "Archiving tour..."
RESPONSE=$(curl -s -X POST "$BASE_URL/tours/$TOUR_ID/archive" \
  -H "Authorization: Bearer $TOKEN")
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"

# --- REACTIVATE TOUR ---
echo "Reactivating tour..."
RESPONSE=$(curl -s -X POST "$BASE_URL/tours/$TOUR_ID/reactivate" \
  -H "Authorization: Bearer $TOKEN")
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"

# --- VIEW TOURS AS AUTHOR ---
echo "Viewing tours as author..."
RESPONSE=$(curl -s -X GET "$BASE_URL/tours/authors/$USER_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
