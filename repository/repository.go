// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

// NewRepository creates a new Repository instance with the provided options.
// It initializes a connection to the PostgreSQL database using the provided
// data source name (DSN) and verifies the connection is successful.
// If any errors occur during initialization, the function will panic, which is normal,
// considering that database is integral part of the application.
// The returned Repository instance contains the initialized database connection.
func NewRepository(opts NewRepositoryOptions) *Repository {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		log.Printf("error init postgres %s", err.Error())
		panic(err)
	} else {
		log.Printf("successfully init postgres")
	}

	if err = db.Ping(); err != nil {
		log.Printf("error ping postgres %s", err.Error())
		panic(err)
	} else {
		log.Printf("successfully ping postgres")
	}

	if db != nil {
		log.Printf("successfully connect to postgres")
	}
	return &Repository{
		Db: db,
	}
}
