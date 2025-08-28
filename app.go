package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func HandleError(err error) {
	if err != nil {
		log.Println("We've encountered an error: ", err)
	}
}

// Create connection string and open a db connection
func (app *App) Initialise(DbUser string, DbPassword string, DbName string) error {

	//Initialise DB connection
	ConnectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", ConnectionString)
	HandleError(err)
	// defer A.DB.Close()

	//Initialise Router

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router)) //Runs until an error
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/", app.homepage)
	app.Router.HandleFunc("/products", app.allProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.oneProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.addProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

// payload interface{} = any type
func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// To handle not ok and not 200 errors
func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

// handlers
func (app *App) homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: Homepage")
	fmt.Fprintf(w, "Welcome Home")
}

func (app *App) allProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: AllProducts")
	products, err := allProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) oneProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid Product ID")
		return
	}
	p := Product{ID: id}
	p.oneProduct(app.DB)
	sendResponse(w, http.StatusOK, p)
}

func (app *App) addProduct(w http.ResponseWriter, r *http.Request) {
	var prod Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	//NewDecoder returns a new decoder that reads from r.
	//Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	if err != nil {
		log.Println("Error: Could not add product")
		sendError(w, http.StatusBadRequest, "Invalid Request Payload") //400
		return
	}

	err = prod.addProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Product %v added!", prod.Name)
	sendResponse(w, http.StatusCreated, prod) //201
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	//get id of the existing product
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid Product ID")
		return
	}

	//get the new json from r
	var prod Product
	err = json.NewDecoder(r.Body).Decode(&prod)
	//NewDecoder returns a new decoder that reads from r.
	//Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	if err != nil {
		log.Println("Error: Could not add product")
		sendError(w, http.StatusBadRequest, "Invalid Request Payload") //400
		return
	}
	prod.ID = id
	err = prod.updateProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Product %v updated!", prod.Name)
	sendResponse(w, http.StatusOK, prod)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid Product ID")
		return
	}
	var prod Product
	prod.ID = id
	err = prod.deleteProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusNotFound, "Invalid Product ID")
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "Successful Deletion"})
}
