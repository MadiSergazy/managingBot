package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"madi_telegram_bot/config"
	"madi_telegram_bot/models"
)

// Deals with database interactions and provides a layer of abstraction. The db.go file includes functions to establish a connection to the MySQL database, execute queries, and retrieve data.
type Database struct {
	connection *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
		return nil, err
	}

	return &Database{
		connection: db,
	}, nil
}

func (db *Database) Close() {
	db.connection.Close()
}

func (db *Database) Execute(query string, args ...interface{}) error {
	res, err := db.connection.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rowsAffected: %v", err)
	}
	if rowsAffected != 1 {
		return errors.New("exec no one query")
	}
	return nil
}

func (db *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.connection.QueryRow(query, args...)
}

func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return rows, nil
}

func (db *Database) InsertIdentifier(identifier int, tableName string) error {
	query := ""
	switch tableName {
	case "admins":
		query = "INSERT INTO admins (identifier) VALUES (?)"
	case "workers":
		query = "INSERT INTO workers (identifier) VALUES (?)"
	default:
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	err := db.Execute(query, identifier)
	if err != nil {
		log.Printf("Error inserting identifier into %s table: %s", tableName, err.Error())
		return err
	}

	log.Printf("Identifier inserted successfully into %s table", tableName)
	return nil
}

// todo: do me right
func InsertLift(dbConnection db.Database, lift models.Lift) (int, error) {
	// Check if the worker exists in the database
	var workerID int
	err := dbConnection.QueryRow("SELECT id FROM workers WHERE phone_number = ?", lift.WorkerPhoneNumber).Scan(&workerID)
	if err != nil {
		return 0, err
	}

	// Insert the lift into the database and return the generated ID
	result, err := dbConnection.Exec("INSERT INTO lifts (name_lift, worker_id) VALUES (?, ?)", lift.Name, workerID)
	if err != nil {
		return 0, err
	}

	liftID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(liftID), nil
}

func InsertTask(dbConnection db.Database, task models.Task) (int, error) {
	// Check if the worker exists in the database
	var workerID int
	err := dbConnection.QueryRow("SELECT id FROM workers WHERE phone_number = ?", task.WorkerPhoneNumber).Scan(&workerID)
	if err != nil {
		return 0, err
	}

	// Insert the task into the database and return the generated ID
	result, err := dbConnection.Exec("INSERT INTO tasks (nameOfTask, dateStart, dateEnd, isDone, worker_id) VALUES (?, ?, ?, ?, ?)", task.Name, task.StartDate, task.EndDate, task.IsDone, workerID)
	if err != nil {
		return 0, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(taskID), nil
}
