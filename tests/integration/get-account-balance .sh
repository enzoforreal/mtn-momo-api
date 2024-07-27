#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"

# Charger le TOKEN extrait
token=$(cat /tmp/momo_auth_token)

# Vérifier si le TOKEN est vide
if [ -z "$token" ]; then
    echo "TOKEN is missing"
    exit 1
fi

# Afficher les variables d'environnement chargées et le TOKEN
echo "TOKEN: $token"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"
echo "ENVIRONMENT: $ENVIRONMENT"

# Exécuter la requête avec des détails de débogage
response=$(curl -s -X GET "http://localhost:8080/get-account-balance" \
    -H "Authorization: Bearer $token" \
    -H "X-Target-Environment: $ENVIRONMENT" \
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY")

# Afficher la réponse
echo "Response: $response"
