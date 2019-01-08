package bank

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nikitsenka/bank-go/bank/utils"
	"log"
	"strconv"
)

type Service interface {
	CreateClient(client Client) (Client, error)
	CreateTransaction(trans Transaction) (Transaction, error)
	GetBalance(clientID int) (Balance, error)
	Close()
}

var DB_HOST = utils.GetEnv("POSTGRES_HOST", "localhost")
var DB_USER = utils.GetEnv("POSTGRES_USER", "postgres")
var DB_PASSWORD = utils.GetEnv("POSTGRES_PASSWORD", "test1234")
var DB_NAME = utils.GetEnv("POSTGRES_NAME", "postgres")
var POOL_SIZE = utils.GetEnv("POSTGRES_CONNECTIONS_SIZE", "95")
var POOL_IDLE_SIZE = utils.GetEnv("POSTGRES_CONNECTIONS_IDLE_SIZE", "2")

type dbClient struct {
	db *sql.DB
}

//New creates postgres implementation of simple bank service
func NewDbService() Service {
	return &dbClient{newDb()}
}

//func Init() {
//	db, _ := newDb()
//	var e error
//	_, e = db.Query("DROP TABLE IF EXISTS client")
//	_, e = db.Query("DROP TABLE IF EXISTS account")
//	_, e = db.Query("DROP TABLE IF EXISTS transaction")
//	checkErr(e)
//	_, e = db.Query("CREATE TABLE client(id SERIAL PRIMARY KEY NOT NULL, name VARCHAR(20), email VARCHAR(20), phone VARCHAR(20));")
//	_, e = db.Query("CREATE TABLE transaction(id SERIAL PRIMARY KEY NOT NULL, from_client_id INTEGER, to_client_id INTEGER, amount INTEGER);")
//	checkErr(e)
//	db.Close()
//}

func (s *dbClient) CreateClient(client Client) (Client, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO client(name, email, phone) VALUES ($1, $2, $3) RETURNING id",
		client.Name, client.Email, client.Phone).Scan(&id)
	fmt.Println("Created client with id", id)
	if err != nil {
		return client, err
	}
	client.Id = id
	return client, nil
}

func (s *dbClient) CreateTransaction(trans Transaction) (Transaction, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO transaction(from_client_id, to_client_id, amount) VALUES ($1, $2, $3) RETURNING id",
		trans.From_client_id, trans.To_client_id, trans.Amount).Scan(&id)
	fmt.Println("Created transaction with id", id)
	if err != nil {
		return trans, err
	}
	trans.Id = id
	return trans, nil
}

func (s *dbClient) GetBalance(client_id int) (Balance, error) {
	var balance Balance
	err := s.db.QueryRow(`
				SELECT debit - credit
				FROM
				  (
					SELECT COALESCE(sum(amount), 0) AS debit
					FROM transaction
					WHERE to_client_id = $1
				  ) a,
				  (
					SELECT COALESCE(sum(amount), 0) AS credit
					FROM transaction
					WHERE from_client_id = $1
				  ) b;
		`, client_id).Scan(&balance.Balance)
	fmt.Println("Calculated balance with client id", client_id)
	return balance, err
}

func (s *dbClient) Close() {
	s.db.Close()
}

func newDb() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
		panic("Error creating postgres conection")
	}
	// Postgres by default has 100 max connection. Sum of max open and max iddle connection should be less than postgres max connection
	poolSize, err := strconv.Atoi(POOL_SIZE)
	if err != nil {
		log.Fatal(err)
		panic("Error getting pool size")
	}
	db.SetMaxOpenConns(poolSize)
	poolIdleSize, err := strconv.Atoi(POOL_IDLE_SIZE)
	if err != nil {
		log.Fatal(err)
		panic("Error getting pool idle size")
	}
	db.SetMaxIdleConns(poolIdleSize)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		panic("Error pinging database")
	}
	return db
}
