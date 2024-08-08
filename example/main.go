package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/enzoforreal/mtn-momo-api/momo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("/home/renaud/mtn-momo-api/.env")
	if err != nil {
		log.Fatalf("Error loading .env file :%v", err)
	}
}

func main() {
	router := gin.Default()

	router.POST("/create-api-user", createAPIUserHandler)
	router.POST("/create-api-key", createAPIKeyHandler)
	router.POST("/get-auth-token", getAuthTokenHandler)
	router.POST("/request-to-pay", requestToPayHandler)
	router.POST("/create-oauth2-token", createOauth2TokenHandler)
	router.GET("/payment-status/:reference_id", getPaymentStatusHandler)
	router.GET("/get-account-balance", getAccountBalanceHandler)

	router.Run(":8080")
}

func createAPIUserHandler(c *gin.Context) {
	client := momo.NewClient()
	var req struct {
		ReferenceID  string `json:"reference_id"`
		CallbackHost string `json:"callback_host"`
	}
	if err := c.BindJSON(&req); err != nil {
		momo.HandleError(c, http.StatusBadRequest, err)
		return
	}

	if req.ReferenceID == "" {
		req.ReferenceID = uuid.New().String()
	}

	log.Printf("Creating API user with reference ID %s and callback host %s", req.ReferenceID, req.CallbackHost)
	if err := client.CreateAPIUser(req.ReferenceID, req.CallbackHost); err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println("API user created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "API user created successfully", "reference_id": req.ReferenceID})
}

func createAPIKeyHandler(c *gin.Context) {
	client := momo.NewClient()
	var req struct {
		ReferenceID string `json:"reference_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		momo.HandleError(c, http.StatusBadRequest, err)
		return
	}

	if req.ReferenceID == "" {
		momo.HandleError(c, http.StatusBadRequest, "Reference ID is required")
		return
	}

	log.Printf("Creating API key for reference ID %s", req.ReferenceID)
	apiKey, err := client.CreateAPIKey(req.ReferenceID)
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println("API key created successfully")
	c.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
}

func getAuthTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		momo.HandleError(c, http.StatusBadRequest, "Authorization header missing")
		return
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		momo.HandleError(c, http.StatusBadRequest, "Invalid authorization header")
		return
	}

	authParts := strings.SplitN(string(decodedAuth), ":", 2)
	if len(authParts) != 2 {
		momo.HandleError(c, http.StatusBadRequest, "Invalid authorization format")
		return
	}

	client := momo.NewClient()
	authToken, err := client.GetAuthToken()
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println("Token retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"token": authToken.AccessToken, "expires_in": authToken.ExpiresIn})
}

func getAccountBalanceHandler(c *gin.Context) {
	client := momo.NewClient()
	token := c.GetHeader("Authorization")
	if token == "" {
		momo.HandleError(c, http.StatusBadRequest, "Authorization header missing")
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	balance, err := client.GetAccountBalance(token)
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println("Account balance retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func requestToPayHandler(c *gin.Context) {
	client := momo.NewClient()
	token := c.GetHeader("Authorization")
	if token == "" {
		momo.HandleError(c, http.StatusBadRequest, "Authorization header missing")
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	var req momo.RequestToPay
	if err := c.BindJSON(&req); err != nil {
		momo.HandleError(c, http.StatusBadRequest, err)
		return
	}

	referenceID, err := client.RequestToPay(token, req)
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Payment request created successfully", "reference_id": referenceID})
}

func createOauth2TokenHandler(c *gin.Context) {
	var req struct {
		AuthReqID string `form:"auth_req_id" binding:"required"`
	}
	if err := c.Bind(&req); err != nil {
		momo.HandleError(c, http.StatusBadRequest, err)
		return
	}

	client := momo.NewClient()
	oauth2Token, err := client.CreateOauth2Token(req.AuthReqID)
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println("OAuth2 token retrieved successfully")
	c.JSON(http.StatusOK, oauth2Token)
}

func getPaymentStatusHandler(c *gin.Context) {
	client := momo.NewClient()
	token := c.GetHeader("Authorization")
	if token == "" {
		momo.HandleError(c, http.StatusBadRequest, "Authorization header missing")
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	referenceID := c.Param("reference_id")
	if referenceID == "" {
		momo.HandleError(c, http.StatusBadRequest, "Reference ID is required")
		return
	}

	paymentStatus, err := client.GetPaymentStatus(referenceID, token)
	if err != nil {
		momo.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Payment status for reference ID %s retrieved successfully", referenceID)
	c.JSON(http.StatusOK, paymentStatus)
}
