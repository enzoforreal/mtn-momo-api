package main

import (
	"fmt"
	"log"

	"github.com/enzoforreal/mtn-momo-api/momo"
)

func main() {
	// Créez un nouveau client avec votre clé API et l'environnement cible
	client := momo.NewClient("your-api-key", "sandbox")

	// Obtenez un token d'authentification
	token, err := client.GetAuthToken()
	if err != nil {
		log.Fatalf("Erreur lors de l'obtention du token d'authentification: %v", err)
	}

	// Affichez le token d'authentification pour confirmation
	fmt.Printf("Token d'authentification: %s\n", token.AccessToken)

	// Obtenez le solde de votre compte
	balance, err := client.GetAccountBalance(token.AccessToken)
	if err != nil {
		log.Fatalf("Erreur lors de l'obtention du solde du compte: %v", err)
	}

	// Affichez le solde disponible
	fmt.Printf("Solde disponible: %s %s\n", balance.AvailableBalance, balance.Currency)

	// Créez une requête de paiement
	request := momo.RequestToPay{
		Amount:     "100",
		Currency:   "USD",
		ExternalId: "7890",
		Payer: momo.Payer{
			PartyIdType: "MSISDN",
			PartyId:     "1234567890",
		},
		PayerMessage: "Paiement pour services",
		PayeeNote:    "Merci pour votre service",
	}

	// Envoyez la requête de paiement
	result, err := client.RequestToPay(token.AccessToken, request)
	if err != nil {
		log.Fatalf("Erreur lors de la demande de paiement: %v", err)
	}

	// Affichez le statut du paiement
	fmt.Printf("Statut du paiement: %s\n", result.Status)
}
