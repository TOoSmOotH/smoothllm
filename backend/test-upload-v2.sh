#!/bin/bash

# Test script for media upload functionality

# Generate a unique email using timestamp
TIMESTAMP=$(date +%s)
EMAIL="testuser${TIMESTAMP}@example.com"
USERNAME="testuser${TIMESTAMP}"

echo "Using email: $EMAIL"
echo "Using username: $USERNAME"

# First, register and login to get a token
echo "Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${EMAIL}\",
    \"password\": \"testpassword123\",
    \"username\": \"${USERNAME}\"
  }")

echo "Register response: $REGISTER_RESPONSE"

# Login to get token
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${EMAIL}\",
    \"password\": \"testpassword123\"
  }")

echo "Login response: $LOGIN_RESPONSE"

# Extract token (register response has token directly, login response has token field)
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//' | head -1)

if [ -z "$TOKEN" ]; then
  echo "Failed to get token"
  exit 1
fi

echo "Token obtained successfully"

# Create a simple test image (1x1 pixel PNG)
echo "Creating test image..."
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==" | base64 -d > test-avatar.png

# Test avatar upload
echo "Testing avatar upload..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/media/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@test-avatar.png")

echo "Upload response: $UPLOAD_RESPONSE"

# Extract media ID from response
MEDIA_ID=$(echo $UPLOAD_RESPONSE | grep -o '"media_id":[0-9]*' | sed 's/"media_id"://')

if [ -n "$MEDIA_ID" ]; then
  echo "✅ Successfully uploaded avatar with media ID: $MEDIA_ID"
  
  # Test getting media info
  echo "Testing get media info..."
  GET_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/media/$MEDIA_ID \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Get media response: $GET_RESPONSE"
  
  # Extract file path from response
  FILE_PATH=$(echo $GET_RESPONSE | grep -o '"file_path":"[^"]*' | sed 's/"file_path":"//')
  echo "File path: $FILE_PATH"
  
  # Test serving media file
  if [ -n "$FILE_PATH" ]; then
    echo "Testing serve media file..."
    SERVE_RESPONSE=$(curl -s -X GET "http://localhost:8080/$FILE_PATH")
    if [ -n "$SERVE_RESPONSE" ]; then
      echo "✅ Successfully served media file (size: $(echo $SERVE_RESPONSE | wc -c) bytes)"
    else
      echo "❌ Failed to serve media file"
    fi
  fi
  
  # Test deleting media
  echo "Testing delete media..."
  DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/v1/media/$MEDIA_ID \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Delete response: $DELETE_RESPONSE"
  
  if echo "$DELETE_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Successfully deleted media"
  else
    echo "❌ Failed to delete media"
  fi
else
  echo "❌ Failed to upload avatar"
  echo "Full upload response: $UPLOAD_RESPONSE"
fi

# Test cover photo upload
echo "Testing cover photo upload..."
COVER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/media/cover \
  -H "Authorization: Bearer $TOKEN" \
  -F "cover=@test-avatar.png")

echo "Cover upload response: $COVER_RESPONSE"

# Clean up
rm -f test-avatar.png

echo "Test completed!"