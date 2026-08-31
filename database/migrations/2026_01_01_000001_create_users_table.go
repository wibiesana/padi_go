package migrations

import (
	"database/sql"

	"github.com/wibiesana/padi_go_core/database"
	"github.com/wibiesana/padi_go_core/migrator"
)

func init() {
	migrator.Register(
		"2026_01_01_000001_create_users_table",
		func(db *sql.DB) error {
			driver := database.GetDriver()
			var createSQL string

			if driver == "postgres" {
				createSQL = `
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					email VARCHAR(255) UNIQUE NOT NULL,
					password VARCHAR(255) NOT NULL,
					role VARCHAR(50) DEFAULT 'user',
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				`
			} else if driver == "mysql" {
				createSQL = `
				CREATE TABLE IF NOT EXISTS users (
					id INT AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					email VARCHAR(255) UNIQUE NOT NULL,
					password VARCHAR(255) NOT NULL,
					role VARCHAR(50) DEFAULT 'user',
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
				`
			} else {
				createSQL = `
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					email TEXT UNIQUE NOT NULL,
					password TEXT NOT NULL,
					role TEXT DEFAULT 'user',
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
				`
			}

			return migrator.ExecHelper(db, createSQL)
		},
		func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS users")
			return err
		},
	)
}
