#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"

# Charger le API_KEY extrait
api_key=$(cat /tmp/momo_api_key)

# Vérifier si le API_KEY est vide
if [ -z "$api_key" ]; then
    echo "API_KEY is missing"
    exit 1
fi

# Afficher les variables d'environnement chargées et le API_KEY
echo "API_KEY: $api_key"
echo "API_USER_ID: $API_USER_ID"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"

# Exécuter la requête avec des détails de débogage
response=$(curl -s -X POST "http://localhost:8080/get-auth-token" \
    -H "Authorization: Basic $(echo -n "$API_USER_ID:$api_key" | base64 | tr -d '\n')" \
    -H "Content-Type: application/json")

# Afficher la réponse
echo "Response: $response"

# Extraire le TOKEN de la réponse
token=$(echo $response | jq -r '.token')

# Vérifier si le TOKEN est vide
if [ -z "$token" ]; then
    echo "Failed to extract TOKEN"
    exit 1
fi

echo "Extracted TOKEN: $token"

# Sauvegarder le TOKEN dans un fichier temporaire
echo $token > /tmp/momo_auth_token
