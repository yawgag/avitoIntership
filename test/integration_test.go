package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"orderPickupPoint/internal/models"
	"os"
	"testing"
	"time"
)

func TestPvzLogic(t *testing.T) {
	baseURL := os.Getenv("SERVER_ADDRESS")
	if baseURL == "" {
		t.Fatal("SERVER_ADDRESS env is not set")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookie jar: %v", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	login_jsonBody := []byte(`{"role":"moderator"}`)
	login_resp, err := client.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(login_jsonBody))
	if err != nil {
		t.Fatalf("failed to send POST request(/dummyLogin): %v", err)
	}
	defer login_resp.Body.Close()

	if login_resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code. got: %d, want: %d", login_resp.StatusCode, http.StatusOK)
	}

	createPvz_jsonBody := []byte(`{"city":"Москва"}`)
	createPvz_resp, err := client.Post(baseURL+"/pvz", "application/json", bytes.NewBuffer(createPvz_jsonBody))
	if err != nil {
		t.Fatalf("failed to send POST request(/pvz): %v", err)
	}
	defer createPvz_resp.Body.Close()

	if createPvz_resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code. got: %d, want: %d", createPvz_resp.StatusCode, http.StatusOK)
	}

	var createPvz_jsonResp models.PickupPointAPI
	if err := json.NewDecoder(createPvz_resp.Body).Decode(&createPvz_jsonResp); err != nil {
		body, _ := io.ReadAll(createPvz_resp.Body)
		t.Fatalf("failed to decode JSON response: %v, raw body: %s", err, string(body))
	}

	createReception_jsonBody := []byte(fmt.Sprintf(`{"pvzId":"%s"}`, createPvz_jsonResp.Id))
	createReception_resp, err := client.Post(baseURL+"/receptions", "application/json", bytes.NewBuffer(createReception_jsonBody))
	if err != nil {
		t.Fatalf("failed to send POST request(/receptions): %v", err)
	}
	defer createReception_resp.Body.Close()

	var createReception_jsonResp models.ReceptionAPI
	if err := json.NewDecoder(createReception_resp.Body).Decode(&createReception_jsonResp); err != nil {
		body, _ := io.ReadAll(createReception_resp.Body)
		t.Fatalf("failed to decode JSON response: %v, raw body: %s", err, string(body))
	}

	numberOfProducts := 10
	productTypes := []string{"электроника", "одежда", "обувь"}
	products := make([]models.ProductAPI, 0, numberOfProducts+1)

	for _ = range numberOfProducts {
		AddProduct_jsonBody := []byte(fmt.Sprintf(`{"pvzId":"%s","type":"%s"}`, createPvz_jsonResp.Id, productTypes[rand.Intn(len(productTypes))]))

		AddProduct_resp, err := client.Post(baseURL+"/products", "application/json", bytes.NewBuffer(AddProduct_jsonBody))
		if err != nil {
			t.Fatalf("failed to send POST request(/products): %v", err)
		}
		defer AddProduct_resp.Body.Close()
		var AddProduct_jsonResp models.ProductAPI

		if err := json.NewDecoder(AddProduct_resp.Body).Decode(&AddProduct_jsonResp); err != nil {
			body, _ := io.ReadAll(AddProduct_resp.Body)
			t.Fatalf("failed to decode JSON response: %v, raw body: %s", err, string(body))
		}
		products = append(products, AddProduct_jsonResp)

	}

	for num, item := range products {
		t.Logf("new item %d: %s", num, item)
	}

	closeReceptionUrl := fmt.Sprintf("/pvz/%s/close_last_reception", createPvz_jsonResp.Id)
	closeReception_resp, err := client.Post(baseURL+closeReceptionUrl, "application/json", nil)
	if err != nil {
		t.Fatalf("failed to send POST request(/receptions): %v", err)
	}
	defer closeReception_resp.Body.Close()
	t.Log("close reception")
}
