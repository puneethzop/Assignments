package migrations

import "gofr.dev/pkg/gofr/migration"

const createTask = `CREATE TABLE IF NOT EXISTS TASKS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		task TEXT,
		completed BOOL DEFAULT FALSE,
		user_id int
	);`

func create_tasks_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTask)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
