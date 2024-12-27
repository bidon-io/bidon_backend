#!/bin/sh

SUPERUSER_LOGIN=${SUPERUSER_LOGIN}
SUPERUSER_PASSWORD=${SUPERUSER_PASSWORD}

ADMIN_EMAIL="${SUPERUSER_LOGIN}@example.com"
ADMIN_PASSWORD=${SUPERUSER_PASSWORD}
API_URL="http://localhost:3200"

if [ -z "$SUPERUSER_LOGIN" ] || [ -z "$SUPERUSER_PASSWORD" ]; then
  echo "Error: SUPERUSER_LOGIN and SUPERUSER_PASSWORD environment variables are required but not set."
  exit 1
fi

AUTH_HEADER=$(echo -n "$SUPERUSER_LOGIN:$SUPERUSER_PASSWORD" | base64)

echo "Waiting for bidon-admin API to be ready..."
while ! curl -s "$API_URL/health_checks" | grep 'OK' > /dev/null; do
  sleep 5
done

echo "Checking if admin user exists..."
ADMIN_EXISTS=$(curl -s -H "Authorization: Basic $AUTH_HEADER" "$API_URL/api/users" | grep "$ADMIN_EMAIL")

if [ -z "$ADMIN_EXISTS" ]; then
  echo "Creating admin user..."
  RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/curl_response -X POST \
    "$API_URL/api/users" \
    -H "Authorization: Basic $AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "'"$ADMIN_EMAIL"'",
      "password": "'"$ADMIN_PASSWORD"'",
      "is_admin": true
    }')
  HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

  if [ "$HTTP_CODE" -eq 201 ]; then
    echo "Admin user created successfully."
  else
    echo "Failed to create admin user. Server responded with HTTP $HTTP_CODE."
    echo "Response body:"
    cat /tmp/curl_response
    exit 1
  fi
else
  echo "Admin user already exists. Skipping initialization."
fi
