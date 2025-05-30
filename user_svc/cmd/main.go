package main

import (
	"fmt"
	"log"

	"github.com/order_management/user_svc/internal/db"
)

func main() {
	fmt.Println("Hello World!")
	_, err := db.InitDB()
	if err != nil {
		log.Fatal("Error to :", err)
	}
}
