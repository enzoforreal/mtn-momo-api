#!/bin/bash

# Remplacez par vos propres valeurs
API_USER_ID="c67802a1-8221-4f51-8bf5-64362c20c34e"
API_KEY="498ff0fcce3240a3981de3d923f887a0"
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
