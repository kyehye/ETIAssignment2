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

//Updating existing tutor information
func UpdateTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("UPDATE Tutors SET Name='%s', Description='%s' WHERE TutorID='%s'",
		t.Name, t.Description, t.TutorID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

/*
//Updating existing tutor information
func ViewTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("SELECT * FROM Tutors WHERE TutorID='%s'",
		t.TutorID, t.Name, t.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}*/

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

/*
//Searches tutor based on personal information
func SearchTutors(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("SELECT * FROM Tutors WHERE TutorID='%s'",
		t.TutorID, t.Name, t.Description)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}


//Tutor using mobile phone number to login
func TutorLogin(db *sql.DB, TutorID string) (Tutor, string) {
	query := fmt.Sprintf("SELECT * FROM Tutor WHERE TutorID = '%s'", TutorID)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&tutors.TutorID, &tutors.FirstName, &tutors.LastName, &tutors.MobileNo, &tutors.EmailAdd); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Tutor login failed."
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
		// POST is for creating new tutor
		if r.Method == "POST" {
			// read the string sent to the service
			var newTutor Tutor
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTutor)
				//Check if user fill up the required information for registering tutor's account
				if newTutor.TutorID == "" || newTutor.Name == "" || newTutor.Description == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply tutor " + "information " + "in JSON format"))
					return
				} else {
					CreateNewTutor(db, newTutor) //Once everything is checked, tutor's account will be created
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Successfully created tutor's account"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply tutor information " +
					"in JSON format"))
			}
		}
		//---PUT is for creating or updating existing tutor---
		if r.Method == "PUT" {
			queryParams := r.URL.Query() //used to resolve the conflict of calling API using the '%s'?TutorID='%s' method
			TutorID = queryParams["TutorID"][0]
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &tutors)
				//Check if user fill up the required information for updating tutor's account information
				if tutors.TutorID == "" || tutors.Name == "" || tutors.Description == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply tutor " + " information " + "in JSON format"))
				} else {
					tutors.TutorID = TutorID
					UpdateTutor(db, tutors)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Successfully updated tutor's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "tutor information " + "in JSON format"))
			}
		}

	}
	//---Deny any deletion of tutor's account or other tutor's information
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("403 - For audit purposes, tutor's account cannot be deleted."))
	}
}

/*
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
*/

//Function main for testing purposes only
func main() {
	// instantiate tutor
	tutors, err := sql.Open("mysql", "user:password@tcp(172.20.30.96:8022)/tutors")
	db = tutors
	// handle error
	if err != nil {
		panic(err.Error())
	}
	//handle the API connection across tutor's microservices
	router := mux.NewRouter()
	router.HandleFunc("/tutors", tutor).Methods(
		"GET", "POST", "PUT", "DELETE")
	/*router.HandleFunc("/tutors/listtutors", listtutor).Methods(
		"GET")
	router.HandleFunc("/tutors/searchtutors", searchtutor).Methods(
		"GET")*/
	fmt.Println("Tutors microservice API --> Listening at port 8011")
	log.Fatal(http.ListenAndServe(":8011", router))

	defer db.Close()
}
