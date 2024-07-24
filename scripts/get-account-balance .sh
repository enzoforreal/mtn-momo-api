#!/bin/bash

# Remplacez par vos propres valeurs
ACCESS_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJSMjU2In0.eyJjbGllbnRJZCI6IjQ2NjgwYzIzLTVjYjgtNGY2ZS04Zjc1LWFlY2FhNmE3ZDQxNSIsImV4cGlyZXMiOiIyMDI0LTA3LTI0VDE1OjQwOjA3LjI1MyIsInNlc3Npb25JZCI6ImJmNDVkMmIxLTMxN2ItNDMwMi05ZjExLTdhMzljNjdiMTM4ZSJ9.LXJT4zoXUOA9e2z5VSZkjGG81jvvHjWnl9AgZXHouum8ImSLd_Te-_lzEuWkDLZoqa5N65-PTc5E5Zo_-cKq67cria6su1ajGhrcReZhjPW3brTW4AZoVdFQmHBT7D_9-NyM_XbuKei0TswcyDZx8gVbenwp6OAnFJ420h74BWzAeusFjYvKUlgaBKY3LVaN6Gd4R5sVyNN465bKHEaMLceaoOAwR4dRUXIxp2wKrKEtUUDXMdYQ5hhP4fce68D7cPbqvcBeLBmkdKQSzj0MVzIGT4nhCjxm9YkFtX3XNqPYthv3ztnMbumkQX8UtfMiamN14GVGIFY27G_10kNPgw"
X_TARGET_ENVIRONMENT="sandbox"
SUBSCRIPTION_KEY="0285a68a2e9542ae8fb41d6512172362"

# Affichage des valeurs pour le débogage
echo "Access Token: $ACCESS_TOKEN"
echo "X-Target-Environment: $X_TARGET_ENVIRONMENT"
echo "Subscription Key: $SUBSCRIPTION_KEY"

# Faites la requête GET pour récupérer le solde du compte
response=$(curl -v -X GET "http://localhost:8080/get-account-balance" \
    -H "Authorization: $ACCESS_TOKEN" \
    -H "X-Target-Environment: $X_TARGET_ENVIRONMENT" \
    -H "Cache-Control: no-cache" \
    -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY")

echo "response: $response"
