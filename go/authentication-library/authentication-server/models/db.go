package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)



func connect(){
	const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = ""
  dbname   = "authentication"
)
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
if err != nil {
  panic(err)
}
	defer db.Close()
}



