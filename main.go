package main

import (
    "fmt"
	"log"
	"net/http"

    _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	
	"github.com/victortanzy123/govtech-assignment-swe/controller"

)


func main() {
	router := mux.NewRouter()
	
	router.HandleFunc("/api/commonstudents", controller.CommonStudents).Methods("GET")
	router.HandleFunc("/api/register", controller.RegisterStudents).Methods("POST")
	router.HandleFunc("/api/suspend", controller.SuspendStudent).Methods("POST")
	router.HandleFunc("/api/retrievefornotifications", controller.RetrieveForNotification).Methods("POST")

	fmt.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}