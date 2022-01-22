package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Global variable for microservice tutor
var tutors Tutor
var TutorID string
var db *sql.DB

type Tutor struct { // map this type to the record in the table
	TutorID     string
	Name        string
	Description string
}

/* For Console
var ukey string
    db.QueryRow("Select LEFT(MD5(rand()),16)").Scan(&ukey)
    ukey = "T"+ukey
*/
////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////							Functions for MySQL Database										////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
//Registering new tutor
func CreateNewTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("INSERT INTO Tutors VALUES ('%s', '%s', '%s')",
		t.TutorID, t.Name, t.Description)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

//Updating existing passenger information
func UpdateTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("UPDATE Tutors SET Name='%s', Description='%s' WHERE TutorID='%s'",
		t.Name, t.Description, t.TutorID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//Updating existing passenger information
func ViewTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("SELECT * FROM Tutors WHERE TutorID='%s'",
		t.TutorID, t.Name, t.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func DeleteTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("DELETE FROM Tutors WHERE TutorID='%s'",
		t.TutorID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
func ListTutors(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("SELECT * FROM Tutors",
		t.TutorID, t.Name, t.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//Searches tutor based on personal information
func SearchTutors(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("SELECT * FROM Tutors WHERE TutorID='%s'",
		t.TutorID, t.Name, t.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

/*
//Tutor using mobile phone number to login
func TutorLogin(db *sql.DB, MobileNo string) (Tutor, string) {
	query := fmt.Sprintf("SELECT * FROM Tutor WHERE MobileNo = '%s'", MobileNo)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&passengers.PassengerID, &passengers.FirstName, &passengers.LastName, &passengers.MobileNo, &passengers.EmailAdd); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Passenger login failed."
	case nil:
	default:
		panic(err.Error())
	}

	return tutors, errMsg
}
*/

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func tutor(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new passenger
		if r.Method == "POST" {
			// read the string sent to the service
			var newTutor Tutor
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTutor)
				//Check if user fill up the required information for registering Passenger's account
				if newTutor.TutorID == "" || newTutor.Name == "" || newTutor.Description == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				} else {
					CreateNewTutor(db, newTutor) //Once everything is checked, passenger's account will be created
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Successfully created passenger's account"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
			}
		}
		//---PUT is for creating or updating existing passenger---
		if r.Method == "PUT" {
			queryParams := r.URL.Query() //used to resolve the conflict of calling API using the '%s'?PassengerID='%s' method
			TutorID = queryParams["TutorID"][0]
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &tutors)
				//Check if user fill up the required information for updating Passenger's account information
				if tutors.TutorID == "" || tutors.Name == "" || tutors.Description == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + " information " + "in JSON format"))
				} else {
					tutors.TutorID = TutorID
					UpdateTutor(db, tutors)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Successfully updated passenger's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "passenger information " + "in JSON format"))
			}
		}

	}
	if r.Method == "GET" {
		if _, ok := tutors[params["TutorID"]]; ok {
			json.NewEncoder(w).Encode(
				tutors[params["TutorID"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No tutor found"))
		}
	}
	//---Deny any deletion of passenger's account or other passenger's information
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("403 - For audit purposes, passenger's account cannot be deleted."))
	}
}

func listtutor(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if _, ok := tutors[params["TutorID"]]; ok {
			json.NewEncoder(w).Encode(
				tutors[params["TutorID"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No tutor found"))
		}
	}
}

func searchtutor(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if _, ok := tutors[params["TutorID"]]; ok {
			json.NewEncoder(w).Encode(
				tutors[params["TutorID"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No tutor found"))
		}

	}
}

//Function main for testing purposes only
func main() {
	// instantiate passengers
	tutors, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/tutors")

	db = tutors
	// handle error
	if err != nil {
		panic(err.Error())
	}
	//handle the API connection across all three microservices, Passengers, Trips and Drivers
	router := mux.NewRouter()
	router.HandleFunc("/tutors", tutor).Methods(
		"GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/tutors/listtutors", listtutor).Methods(
		"GET")
	router.HandleFunc("/tutors/searchtutors", searchtutor).Methods(
		"GET")
	fmt.Println("Tutors microservice API --> Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

	defer db.Close()
}
