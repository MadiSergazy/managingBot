package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

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

func (db *Database) ExecuteWithLastInsertID(query string, args ...interface{}) (int, error) {
	res, err := db.connection.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rowsAffected: %v", err)
	}
	if rowsAffected != 1 {
		return 0, errors.New("exec no one query")
	}
	rowsAffected, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get rowsAffected: %v", err)
	}
	return int(rowsAffected), nil
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

func (db *Database) InsertIdentifier(identifier int, tableName string, phone string, userName string) error {
	query := ""
	switch tableName {
	case "admins":
		query = "UPDATE admins SET identifier = ?, name = ? WHERE phone_number = ?"
	case "workers":
		query = "UPDATE workers SET identifier = ?, name = ? WHERE phone_number = ?"
	case "hr_manager":
		query = "UPDATE hr_manager SET identifier = ?, name = ? WHERE phone_number = ?"
	default:
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	err := db.Execute(query, identifier, userName, phone)
	if err != nil {
		log.Printf("Error inserting identifier into %s table: %s", tableName, err.Error())
		return err
	}

	log.Printf("Identifier inserted successfully into %s table", tableName)
	return nil
}

// Step 1: Get employeeID from the database using employeePhoneNumber
func (db *Database) GetEmployeeID(phoneNumber string) (int, error) {
	// Query the database to check if the employee already exists
	query := "SELECT id FROM workers WHERE phone_number = ?"
	row := db.QueryRow(query, phoneNumber)

	var employeeID int64
	err := row.Scan(&employeeID)
	if err != nil {
		// Employee exists, return the ID
		if err != sql.ErrNoRows {

			// Employee doesn't exist, insert a new record
			query = "INSERT INTO workers (phone_number) VALUES (?)"
			result, err := db.connection.Exec(query, phoneNumber)
			if err != nil {
				return 0, err
			}

			// Get the ID of the newly inserted employee
			employeeID, err = result.LastInsertId()
			if err != nil {
				return 0, err
			}

			return int(employeeID), nil
		} else {
			return 0, err
		}
	}
	return int(employeeID), nil
}

// Step 2: Insert the lift information into the lifts table
func (db *Database) InsertLiftInfo(employeeID int, liftName string) (int, error) {
	query := "INSERT INTO lifts (name_lift, worker_id) VALUES (?, ?)"
	result, err := db.connection.Exec(query, liftName, employeeID)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted lift
	liftID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(liftID), nil
}

// Step 2a: Insert the project information into the projects table

// Step 2a: Insert the project information into the projects table
func (db *Database) InsertProjectInfo(nameResident string, workerID, liftID int) (int, error) {
	query := "INSERT INTO projects (name_resident, worker_id, lift_id) VALUES (?, ?, ?)"
	result, err := db.connection.Exec(query, nameResident, workerID, liftID)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted project
	projectID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	exec, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if exec != 1 {
		return 0, errors.New("Error inserting project: no one exec")
	}

	return int(projectID), nil
}

// Step 2b: Insert the lift details into the lift_details table
func (db *Database) InsertLiftDetails(nameResident, nameLift string) (int, error) {
	query := "INSERT INTO lift_details (name_resident, name_lift) VALUES (?, ?)"
	res, err := db.connection.Exec(query, nameResident, nameLift)
	if err != nil {
		return 0, err
	}
	lift_details_id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lift_details_id), nil
}

// Step 3: Insert the tasks into the tasks table
func (db *Database) InsertTasks(taskIDToName map[int]string, employeeID int, tasks map[int]models.TaskDates, liftID int) error {
	query := "INSERT INTO tasks (nameOfTask, dateStart, dateEnd, isDone, lift_id) VALUES (?, ?, ?, ?, ?)"
	for taskID, dates := range tasks {
		taskName := taskIDToName[taskID]
		startDate, _ := time.Parse("02/01/2006", dates.StartDate)
		endDate, _ := time.Parse("02/01/2006", dates.EndDate)

		_, err := db.connection.Exec(query, taskName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), 0, liftID)
		if err != nil {
			fmt.Println("Start Date: ", startDate)
			fmt.Println("Finish date: ", endDate)
			fmt.Println("taskIDToName : ", taskIDToName)
			return err
		}
	}
	return nil
}

// Function to view tasks for a specific employee number
func ViewTasksForAdmin(dbConnection Database, employeeNumber string) ([]models.TaskInfo, error) {
	// Query the database to fetch tasks for the specified employee number
	//todo do it right
	query := `
  SELECT
  p.name_resident,
  l.name_lift,
  w.phone_number,
  t.nameOfTask,
  t.dateStart,
  t.dateEnd,
  t.isDone,
  t.id, 
  t.is_validate
FROM projects p
JOIN lifts l ON p.lift_id = l.id
JOIN workers w ON p.worker_id = w.id
JOIN tasks t ON t.lift_id = l.id
WHERE w.phone_number = ?
ORDER BY p.id, t.id;
`

	rows, err := dbConnection.Query(query, employeeNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to store the task information
	taskInfoList := make([]models.TaskInfo, 0)

	// Iterate through the rows and extract task information
	for rows.Next() {
		var (
			elevatorName        string
			residentialComplex  string
			employeePhoneNumber string
			taskName            string
			startDateStr        []uint8
			endDateStr          []uint8
			isDone              bool
			taskID              int
			is_validate         bool
		)

		err := rows.Scan(
			&residentialComplex,
			&elevatorName,
			&employeePhoneNumber,
			&taskName,
			&startDateStr,
			&endDateStr,
			&isDone,
			&taskID,
			&is_validate,
		)
		if err != nil {
			fmt.Println("Scan err: ", err)
			return nil, err
		}

		// Parse the date strings into time.Time objects
		startDate, _ := time.Parse("2006-01-02", string(startDateStr))
		endDate, _ := time.Parse("2006-01-02", string(endDateStr))

		// Create a TaskInfo struct and populate it with the retrieved values
		taskInfo := models.TaskInfo{
			ElevatorName:        elevatorName,
			ResidentialComplex:  residentialComplex,
			EmployeePhoneNumber: employeePhoneNumber,
			TaskName:            taskName,
			StartDate:           startDate.Format("02/01/2006"),
			EndDate:             endDate.Format("02/01/2006"),
			IsDone:              isDone,
			TaskID:              taskID,
			Is_validate:         is_validate,
		}

		taskInfoList = append(taskInfoList, taskInfo)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("rows.Err(): ", err)
		return nil, err
	}

	return taskInfoList, nil
}
