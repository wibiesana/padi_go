package migrations

import (
	"database/sql"

	"github.com/wibiesana/padi_go_core/database"
	"github.com/wibiesana/padi_go_core/migrator"
)

func init() {
	migrator.Register(
		"2026_01_01_000002_create_password_resets_table",
		func(db *sql.DB) error {
			driver := database.GetDriver()
			var createSQL string

			if driver == "postgres" {
				createSQL = `
				CREATE TABLE IF NOT EXISTS password_resets (
					id SERIAL PRIMARY KEY,
					email VARCHAR(255) NOT NULL,
					token VARCHAR(255) NOT NULL,
					expires_at BIGINT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
				CREATE INDEX IF NOT EXISTS idx_pwd_resets_email_token ON password_resets(email, token);
				`
			} else if driver == "mysql" {
				createSQL = `
				CREATE TABLE IF NOT EXISTS password_resets (
					id INT AUTO_INCREMENT PRIMARY KEY,
					email VARCHAR(255) NOT NULL,
					token VARCHAR(255) NOT NULL,
					expires_at BIGINT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					INDEX idx_pwd_resets_email_token (email, token)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
				`
			} else {
				createSQL = `
				CREATE TABLE IF NOT EXISTS password_resets (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					email TEXT NOT NULL,
					token TEXT NOT NULL,
					expires_at INTEGER NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
				CREATE INDEX IF NOT EXISTS idx_pwd_resets_email_token ON password_resets(email, token);
				`
			}

			return migrator.ExecHelper(db, createSQL)
		},
		func(db *sql.DB) error {
			_, err := db.Exec("DROP TABLE IF EXISTS password_resets")
			return err
		},
	)
}
