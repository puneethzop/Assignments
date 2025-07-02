package migrations

import "gofr.dev/pkg/gofr/migration"

const createUser = `REATE TABLE IF NOT EXISTS USERS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100)
	);`

func create_user_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createUser)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
