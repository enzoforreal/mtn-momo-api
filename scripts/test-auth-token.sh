#!/bin/bash

# Remplacez par vos propres valeurs
API_USER_ID="793a623f-71bd-4db1-aec9-02beda09874e"
API_KEY="8ca0221196b44319b2e04b20ab216e32"
SUBSCRIPTION_KEY="0285a68a2e9542ae8fb41d6512172362"

AUTH_STRING="${API_USER_ID}:${API_KEY}"

# Encodage Base64
ENCODED_AUTH=$(echo -n $AUTH_STRING | base64 | tr -d '\n')

# Afficher l'encodage pour vérification
echo "Encoded Auth: $ENCODED_AUTH"

# Faire l'appel API
response=$(curl -v -X POST "http://localhost:8080/get-auth-token" \
    -H "Authorization: Basic $ENCODED_AUTH" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    -H "Content-Type: application/json")

# Afficher la réponse
echo "Response: $response"
