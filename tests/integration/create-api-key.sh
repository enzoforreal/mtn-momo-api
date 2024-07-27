#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"

# Charger le REFERENCE_ID extrait
reference_id=$(cat /tmp/momo_reference_id)

# Vérifier si le REFERENCE_ID est vide
if [ -z "$reference_id" ]; then
    echo "REFERENCE_ID is missing"
    exit 1
fi

# Afficher les variables d'environnement chargées et le REFERENCE_ID
echo "REFERENCE_ID: $reference_id"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"

# Exécuter la requête avec des détails de débogage
response=$(curl -s -X POST "http://localhost:8080/create-api-key" \
    -H "Content-Type: application/json" \
    -H "Cache-Control: no-cache" \
    -d "{
          \"reference_id\": \"$reference_id\"
        }")

# Afficher la réponse
echo "Response: $response"

# Extraire le API_KEY de la réponse
api_key=$(echo $response | jq -r '.api_key')

# Vérifier si le API_KEY est vide
if [ -z "$api_key" ]; then
    echo "Failed to extract API_KEY"
    exit 1
fi

echo "Extracted API_KEY: $api_key"

# Sauvegarder le API_KEY dans un fichier temporaire
echo $api_key > /tmp/momo_api_key
