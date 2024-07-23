#!/bin/bash

# Remplacez par vos propres valeurs
REFERENCE_ID="c67802a1-8221-4f51-8bf5-64362c20c34e"

response=$(curl -s -X POST "http://localhost:8080/create-api-key" \
    -H "Content-Type: application/json" \
    -d "{
          \"reference_id\": \"$REFERENCE_ID\"
        }")

# Afficher la réponse
echo "Response: $response"
