package controllers

import (
	"fmt"
	"log"

	"github.com/markstanden/authentication/models"
)

// Connect is a test function to check connection
  func Connect(email string) *models.User{
    const (
    host     = "localhost"
    port     = 5432
    username = "postgres"
    password = ""
    databaseName   = "authentication"
  )

    db := models.NewConnection(host, username, password, databaseName, port)
    err := db.DB.Ping();
    if err != nil {
      fmt.Println("Connection Failure", err)
    }
    user, err := db.FindByEmail(email)
    if err != nil {
      log.Println("User Not Found")
    }
    return user
  }