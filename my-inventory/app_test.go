package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occurred while initializing Database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
        id int NOT NULL AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        quantity int,
        price float(10,7),
        PRIMARY KEY (id)
    );`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
	log.Println("ClearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v',%v,%v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name": "chair", "quantity": 1, "price": 100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, Got: %v", "chair", m["name"])
	}

	if m["quantity"] != float64(1) {
		t.Errorf("Expected quantity: %v, Got: %v", float64(1), m["name"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	createTable()
	addProduct("chair", 10, 100)

	req, _ := http.NewRequest("GET", "/product/2", nil)
	response := sendRequest(req)
	fmt.Println("GET /product/1 Response:", response.Body.String()) // Debug

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)
	fmt.Println("Parsed oldValue:", oldValue) // Debug

	var product = []byte(`{"name": "connector", "quantity": 10, "price": 100}`)
	req, _ = http.NewRequest("PUT", "/product/2", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)
	fmt.Println("PUT /product/1 Response:", response.Body.String()) // Debug

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)
	fmt.Println("Parsed newValue:", newValue) // Debug

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id: %v, Got: %v", oldValue["id"], newValue["id"])
	}

	if oldValue["name"] == newValue["name"] {
		t.Errorf("Expected name: %v, Got: %v", oldValue["name"], newValue["name"])
	}

	if oldValue["quantity"] != newValue["quantity"] {
		t.Errorf("Expected quantity: %v, Got: %v", oldValue["quantity"], newValue["quantity"])
	}

	if oldValue["price"] != newValue["price"] {
		t.Errorf("Expected price: %v, Got: %v", oldValue["price"], newValue["price"])
	}
}
