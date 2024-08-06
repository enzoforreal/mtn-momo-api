#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

echo "Chargement du fichier d'environnement : $env_file"
if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"
echo "Fichier d'environnement chargé."

# Vérification des variables d'environnement
if [ -z "$SUBSCRIPTION_KEY" ]; then
    echo "SUBSCRIPTION_KEY is missing in environment variables"
    exit 1
fi

# Charger le API_USER_ID extrait
echo "Chargement de API_USER_ID à partir de /tmp/momo_reference_id"
api_user_id=$(cat /tmp/momo_reference_id)
echo "API_USER_ID: $api_user_id"

# Vérifier si le API_USER_ID est vide
if [ -z "$api_user_id" ]; then
    echo "API_USER_ID is missing"
    exit 1
fi

# Charger le API_KEY extrait
echo "Chargement de API_KEY à partir de /tmp/momo_api_key"
api_key=$(cat /tmp/momo_api_key)
echo "API_KEY: $api_key"

# Vérifier si le API_KEY est vide
if [ -z "$api_key" ]; then
    echo "API_KEY is missing"
    exit 1
fi

# Afficher les variables d'environnement chargées, le API_USER_ID et le API_KEY
echo "Variables d'environnement chargées:"
echo "API_KEY: $api_key"
echo "API_USER_ID: $api_user_id"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"

# Vérifier si auth_req_id est passé en argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 auth_req_id"
    exit 1
fi

auth_req_id=$1

# Afficher auth_req_id pour vérification
echo "auth_req_id reçu: $auth_req_id"  # Message de débogage

# Exécuter la requête avec des détails de débogage
response=$(curl -v -X POST "http://localhost:8080/create-oauth2-token" \
    -H "Authorization: Basic $(echo -n "$api_user_id:$api_key" | base64 | tr -d '\n')" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H "X-Target-Environment: sandbox" \
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    --data-urlencode "auth_req_id=$auth_req_id" \
    --data-urlencode "grant_type=urn:openid:params:grant-type:ciba")

# Afficher la réponse
echo "Response: $response"
