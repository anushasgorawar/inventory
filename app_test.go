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

// Released with go 1.4, not to test the main function
func TestMain(m *testing.M) {
	err := a.Initialise("root", "Anusha#13", "testinventory")
	if err != nil {
		log.Fatal("error occured while initialising the database")
	}
	createtable()
	m.Run() //only m method //runs all otehr tests within that package

}

func createtable() {
	query := `CREATE TABLE IF NOT EXISTS products (
	id int NOT NULL AUTO_INCREMENT, 
	name varchar (255) NOT NULL, 
	quantity int, 
	price float (10,7),
	PRIMARY KEY (id)
	)`
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
func clearTable() {
	_, err := a.DB.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal(err)
	}
	_, err = a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	if err != nil {
		log.Fatal(err)
	}
}

func addToTable(name string, quantity int, price float64) {
	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)", name, quantity, price)
	a.DB.Exec(query)
}

func TestGetProduct(t *testing.T) {
	//test responce from get
	clearTable()
	addToTable("guitar", 1, 25.99)
	request, _ := http.NewRequest("GET", "/product/1", nil) //nil is the payload
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

// helper function
func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	//return recorder object and send the rquest in

	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkStatusCode(t *testing.T, expectedStatusCode, recievedStatusCode int) {
	if expectedStatusCode != recievedStatusCode {
		t.Errorf("Expected Status: %v, recieved: %v", expectedStatusCode, recievedStatusCode)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	prod := []byte(`{"Name":"Book","Quantity":2,"Price":0.49}`)
	//body should be of bytebuffer
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(prod))
	request.Header.Set("Content-Type", "application/json")

	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["Name"] != "Book" {
		t.Errorf("Expected name: %v, Got: %v", "Book", m["Name"])
	}
	if m["Quantity"] != 2.0 {
		t.Errorf("Expected quantity: %v, Got: %v", 2.0, m["Quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addToTable("chair", 1, 25.99)

	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addToTable("chair", 1, 25.99)
	request, _ := http.NewRequest("GET", "/product/1", nil) //nil is the payload
	response := sendRequest(request)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name": "connector", "quantity":1, "price":100}`)
	request, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")
	response = sendRequest(request)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["ID"] != newValue["ID"] {
		t.Errorf("Expected quantity: %v, Got: %v", newValue["ID"], oldValue["ID"]) //doesnt change
	}
	if oldValue["Name"] == newValue["Name"] {
		t.Errorf("Expected quantity: %v, Got: %v", newValue["Name"], oldValue["Name"]) //name is updated
	}
}
