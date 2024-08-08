#!/bin/bash

# Charger les variables d'environnement
ENV_FILE="$(dirname "$0")/../../integration.env"
if [ ! -f "$ENV_FILE" ]; then
    echo "Erreur: le fichier $ENV_FILE n'existe pas."
    exit 1
fi
source "$ENV_FILE"

# Vérifier la disponibilité de jq
if ! command -v jq &> /dev/null; then
    echo "Erreur: jq n'est pas installé. Veuillez l'installer pour continuer."
    exit 1
fi

# Fonction pour obtenir le jeton d'authentification
get_auth_token() {
    ./tests/integration/get-auth-token.sh
    if [ ! -f /tmp/momo_auth_token ]; then
        echo "Erreur: impossible d'obtenir le jeton d'authentification"
        exit 1
    fi
}

# Fonction pour créer un utilisateur API
create_api_user() {
    ./tests/integration/create-api-user.sh
    if [ ! -f /tmp/momo_reference_id ]; then
        echo "Erreur: impossible de créer l'utilisateur API"
        exit 1
    fi
}

# Vérifier si le fichier de jeton d'authentification existe, sinon obtenir un nouveau jeton
if [ ! -f /tmp/momo_auth_token ]; then
    echo "Le fichier de jeton d'authentification n'existe pas. Obtention d'un nouveau jeton..."
    get_auth_token
fi

token=$(cat /tmp/momo_auth_token)
echo "Token: $token"

# Vérifier si le jeton est vide
if [ -z "$token" ]; then
  echo "Erreur: le jeton d'authentification est vide"
  exit 1
fi

# Vérifier si le fichier de reference ID existe, sinon créer un nouvel utilisateur API
if [ ! -f /tmp/momo_reference_id ]; then
    echo "Le fichier de reference ID n'existe pas. Création d'un nouvel utilisateur API..."
    create_api_user
fi

reference_id=$(cat /tmp/momo_reference_id)
echo "Reference ID: $reference_id"

# Afficher les valeurs d'environnement pour le débogage
echo "BASE_URL: $BASE_URL"
echo "Environment: $ENVIRONMENT"
echo "Subscription Key: $SUBSCRIPTION_KEY"

# Créer la requête de paiement
response=$(curl -s -w "\nHTTP_STATUS_CODE:%{http_code}\n" -X POST "$BASE_URL/request-to-pay" \
    -H "Authorization: Bearer $token" \
    -H "X-Reference-Id: $reference_id" \
    -H "X-Target-Environment: $ENVIRONMENT" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
    -H "Content-Type: application/json" \
    -H "Cache-Control: no-cache" \
    -d '{
          "amount": "100",
          "currency": "EUR",
          "externalId": "123456",
          "payer": {
            "partyIdType": "MSISDN",
            "partyId": "46733123453"
          },
          "payerMessage": "Payment for invoice 123456",
          "payeeNote": "Invoice 123456 payment"
        }')

# Extraire le corps de la réponse et le code de statut HTTP
response_body=$(echo "$response" | sed -n '1h;1!H;${g;s/\nHTTP_STATUS_CODE:.*//p;}')
http_status=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTP_STATUS_CODE://')

echo "Response body: $response_body"
echo "HTTP status: $http_status"

# Extraire le message de la réponse
message=$(echo "$response_body" | jq -r '.message' 2>/dev/null)

if [ "$http_status" -ne 202 ]; then
  echo "Erreur: la requête de paiement a échoué avec le statut HTTP $http_status"
  exit 1
fi

# Vérifier si le message de succès est reçu
if [ "$message" != "Payment request created successfully" ]; then
  echo "Erreur: la requête de paiement a échoué avec le message $message"
  exit 1
else
  echo "Requête de paiement créée avec succès"

  # Extraire 'reference_id' de la réponse JSON
  reference_id=$(echo "$response_body" | jq -r '.reference_id')

  # Vérifier et afficher le reference_id extrait
  if [ -n "$reference_id" ]; then
      echo "Reference ID extrait: $reference_id"
      echo $reference_id > /tmp/X-Reference-Id-requesttopay
  else
      echo "Erreur: Aucun Reference ID trouvé dans la réponse"
      exit 1
  fi
fi

# Vérifier si le fichier a été créé avec succès
if [ -f /tmp/X-Reference-Id-requesttopay ]; then
    echo "Fichier X-Reference-Id-requesttopay créé avec succès."
else
    echo "Erreur: le fichier X-Reference-Id-requesttopay n'a pas été créé."
    exit 1
fi

# Stocker à nouveau le token au cas où il serait utilisé plus tard
echo $token > /tmp/momo_auth_token
