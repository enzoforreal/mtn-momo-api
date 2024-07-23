#!/bin/bash

# Remplacez par vos propres valeurs
REFERENCE_ID="fa364ea7-cb35-4fe3-8b42-009c22d04fbf"
SUBSCRIPTION_KEY="0285a68a2e9542ae8fb41d6512172362"
CALLBACK_HOST="https://2661-102-141-52-18.ngrok-free.app"

# Exécuter la requête avec des détails de débogage
response=$(curl -v -X POST "http://localhost:8080/create-api-user" \
    -H "X-Reference-Id: $REFERENCE_ID" \
    -H "Content-Type: application/json" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    -d "{
          \"callback_host\": \"$CALLBACK_HOST\"
        }")

# Afficher la réponse
echo "Response: $response"
