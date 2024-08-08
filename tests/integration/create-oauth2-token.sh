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

# Générer un UUID si non fourni
if [ -z "$1" ]; then
    auth_req_id=$(uuidgen)  # Générer un nouveau UUID
    echo "Généré un nouvel auth_req_id: $auth_req_id"
else
    auth_req_id=$1
    echo "auth_req_id reçu: $auth_req_id"
fi

# Afficher les variables d'environnement chargées
echo "Variables d'environnement chargées:"
echo "API_KEY: $api_key"
echo "API_USER_ID: $api_user_id"
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"
echo "auth_req_id: $auth_req_id"

# Exécuter la requête avec des détails de débogage
response=$(curl -s -X POST "http://localhost:8080/create-oauth2-token" \
    -H "Authorization: Basic $(echo -n "$api_user_id:$api_key" | base64 | tr -d '\n')" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H "X-Target-Environment: sandbox" \
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    --data-urlencode "auth_req_id=$auth_req_id" \
    --data-urlencode "grant_type=urn:openid:params:grant-type:ciba")

# Afficher la réponse
echo "Response: $response"
