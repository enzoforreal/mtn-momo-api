#!/bin/bash
env_file="$(dirname "$0")/../../integration.env"

# Vérifier si le fichier d'environnement existe
if [ ! -f "$env_file" ]; then
  echo "Le fichier $env_file est manquant"
  exit 1
fi

source "$env_file"

# Charger le X-Reference-Id-requesttopay
x_reference_id=$(cat /tmp/X-Reference-Id-requesttopay)

# Charger le token d'authentification
token=$(cat /tmp/momo_auth_token)

# Vérifier si le X-Reference-Id-requesttopay est vide
if [ -z "$x_reference_id" ]; then
    echo "X-Reference-Id-requesttopay is missing"
    exit 1
fi

# Afficher les variables d'environnement et les valeurs chargées
echo "SUBSCRIPTION_KEY: $SUBSCRIPTION_KEY"
echo "X-Reference-Id-requesttopay: $x_reference_id"
echo "momo_auth_token: $token"
echo "ENVIRONMENT: $ENVIRONMENT"

# Effectuer la requête cURL pour obtenir le statut du paiement
response=$(curl -s -X GET "http://localhost:8080/payment-status/$x_reference_id" \
            -H "Authorization: Bearer $token" \
            -H "X-Target-Environment: $ENVIRONMENT" \
            -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY")

# Afficher la réponse
echo "Response: $response"

# Inspection de la réponse pour extraire des valeurs
reference_id=$(echo "$response" | jq -r '.referenceId // empty')
status=$(echo "$response" | jq -r '.status // empty')


echo "Extracted Reference ID: $reference_id"
echo "Extracted Status: $status"


# Gestion des erreurs si des valeurs sont vides
if [[ -z "$reference_id" || -z "$status" ]]; then
    echo "Erreur: Données de paiement incomplètes."
    exit 1
fi
