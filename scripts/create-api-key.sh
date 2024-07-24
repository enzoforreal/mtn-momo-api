#!/bin/bash

# Remplacez par vos propres valeurs
REFERENCE_ID="b474a7fc-0446-47c2-b2eb-b3d1fe96daf0"

response=$(curl -v -X POST "http://localhost:8080/create-api-key" \
    -H "Content-Type: application/json" \
    -H "Cache-Control: no-cache" \
    -d "{
          \"reference_id\": \"$REFERENCE_ID\"
        }")


echo "Response: $response"
