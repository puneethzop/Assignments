package datasource

import (
	"database/sql"
)

func New(creds string) (*sql.DB, error) {
	db, err := sql.Open("mysql", creds)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	// Ensure TASKS table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TASKS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		task TEXT,
		completed BOOL DEFAULT FALSE,
		user_id int
	);`)
	if err != nil {
		return nil, err
	}

	// Ensure USERS table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USERS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100)
	);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
