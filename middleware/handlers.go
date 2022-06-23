package middleware

import (
	"database/sql" // package for working with sql
	"encoding/json" // package for encoding and decoding json into struct and vice versa
	"fmt"
	"go-postgres-stocks/models" // models package where Stock schema is defined
	"log"
	"net/http" // used to access the request and response object of the API
	"os"       
	"strconv"
	"github.com/gorilla/mux"    // used to get the params from the API Route
	"github.com/joho/godotenv"  // package used to read the .env file
	_"github.com/lib/pq"        // postgres Golang driver
)

// Response
type response struct {
	ID int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Function for connection with Database
func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// Check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Call to Database!\n")

	return db
}

/* HANDLERS */

// CreateStock create a stock in the Database
func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertStock(stock)

	// Format response message
	res := response{
		ID: insertID,
		Message: "Stock created successfully",
	}

	// Send the response
	json.NewEncoder(w).Encode(res)
}

// GetStock will return a single stock by its id
func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}

	// Send the response
	json.NewEncoder(w).Encode(stock)
}

// GetAllStock will return all the stocks
func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get stocks. %v", err)
	}

	// Send the response
	json.NewEncoder(w).Encode(stocks)
}

// UpdateStock update stock's details in the Database
func UpdateStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", updatedRows)

	// Format response message
	res := response{
		ID: int64(id),
		Message: msg,
	}

	// Send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteStock delete stock's detail in the Database
func DeleteStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	deletedRows := deleteStock(int64(id))

	message := fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", deletedRows)

	// Format the reponse
	res := response{
		ID: int64(id),
		Message: message,
	}

	// Send the response
	json.NewEncoder(w).Encode(res)
}

/* DB Functions */

// Insert one stock in the DB
func insertStock(stock models.Stock) int64 {

	// Connection
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	// Execute SQL
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

// Get one stock from the DB by its stockid
func getStock(id int64) (models.Stock, error) {

	// Connection
	db := createConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	// Execute SQL
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return stock, err
}

// Get all stocks from the DB
func getAllStocks() ([]models.Stock, error) {

	// Connection
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`

	// Execute SQL
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// Close SQL statement
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var stock models.Stock

		// Unmarshal the row object to stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// Append the stock in the stocks slice
		stocks = append(stocks, stock)

	}

	return stocks, err
}

// Update stock (details) in the DB
func updateStock(id int64, stock models.Stock) int64 {

	// Connection
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	// Execute SQL
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// Delete stock in the DB
func deleteStock(id int64) int64 {

	// Connection
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	// Execute SQL
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}