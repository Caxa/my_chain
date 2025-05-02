package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	RegisterRoutes(r)
	return r
}

func TestCreateWalletHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wallet/create", nil)
	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response WalletResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}
	if response.Address == "" || response.PrivateKey == "" {
		t.Error("Expected non-empty wallet response")
	}
}

func TestSendTransactionHandler(t *testing.T) {
	privateKey := "83dab171ec53b5a31c6ddc6f5464e7af0d570290f7c1ec0bc516e91d8daaca0a"

	payload := SendTxRequest{
		PrivateKey: privateKey,
		ToAddress:  "0x1a6Fcb6A2Ae017eC74629085f4D12B38594E6189", // <-- используем настоящий Sepolia-адрес
		Amount:     "0.0001",                                     // ETH
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/eth/send", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := response["txHash"]; !ok {
		t.Error("Expected txHash in response")
	}
}
