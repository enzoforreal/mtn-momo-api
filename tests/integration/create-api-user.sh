#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"

# Afficher les variables d'environnement chargées
echo "REFERENCE_ID: $REFERENCE_ID"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"
echo "CALLBACK_HOST: $CALLBACK_HOST"

# Exécuter la requête avec des détails de débogage
response=$(curl -s -X POST "http://localhost:8080/create-api-user" \
    -H "X-Reference-Id: $REFERENCE_ID" \
    -H "Content-Type: application/json" \
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    -d "{
          \"callback_host\": \"$CALLBACK_HOST\"
        }")

# Afficher la réponse
echo "Response: $response"

# Extraire le REFERENCE_ID de la réponse
reference_id=$(echo $response | jq -r '.reference_id')

# Vérifier si le REFERENCE_ID est vide
if [ -z "$reference_id" ]; then
    echo "Failed to extract REFERENCE_ID"
    exit 1
fi

echo "Extracted REFERENCE_ID: $reference_id"

# Sauvegarder le REFERENCE_ID dans un fichier temporaire
echo $reference_id > /tmp/momo_reference_id
