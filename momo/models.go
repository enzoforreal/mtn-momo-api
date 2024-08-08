package momo

type Client struct {
	ApiKey          string
	ApiUserID       string
	SubscriptionKey string
	Environment     string
}

// Structure pour le token d'authentification
type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Structure pour l'Oauth2Token
type Oauth2TokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int    `json:"expires_in"`
	Scope                 string `json:"scope"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredIn int    `json:"refresh_token_expired_in"`
}

// Structure pour les informations de balance
type Balance struct {
	AvailableBalance string `json:"availableBalance"`
	Currency         string `json:"currency"`
}

// Structure pour les erreurs
type ErrorReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Structure pour une requête de paiement
type RequestToPay struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalId   string `json:"externalId"`
	Payer        Payer  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

type Payer struct {
	PartyIdType string `json:"partyIdType"`
	PartyId     string `json:"partyId"`
}

// Structure pour les résultats de paiement
type RequestToPayResult struct {
	Amount                 string      `json:"amount"`
	Currency               string      `json:"currency"`
	FinancialTransactionId string      `json:"financialTransactionId"`
	ExternalId             string      `json:"externalId"`
	Payer                  Payer       `json:"payer"`
	Status                 string      `json:"status"`
	Reason                 ErrorReason `json:"reason,omitempty"`
}
