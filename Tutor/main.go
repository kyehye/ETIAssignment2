package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
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
	Password    string
}

/* For Console
var TutorID string
db.QueryRow("Select LEFT(MD5(rand()),16)").Scan(&TutorID)
TutorID = "T"+ TutorID
*/
////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////							Functions for MySQL Database										////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
//Registering new tutor
func CreateNewTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("INSERT INTO Tutor VALUES ('%s', '%s', '%s', '%s')",
		t.TutorID, t.Name, t.Description, t.Password)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

//Updating existing tutor information
func UpdateTutor(db *sql.DB, t Tutor) {
	query := fmt.Sprintf("UPDATE Tutor SET Name='%s', Description='%s', Password='%s' WHERE TutorID='%s'",
		t.Name, t.Description, t.Password, t.TutorID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//Delete existing tutor information
func DeleteTutor(db *sql.DB, TutorID string) (Tutor, string) {
	query := fmt.Sprintf("DELETE FROM Tutor WHERE TutorID = '%s'", TutorID)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&tutors.TutorID, &tutors.Name, &tutors.Description, &tutors.Password); err {
	case sql.ErrNoRows:
		errMsg = "TutorID not found."
	case nil:
	default:
		panic(err.Error())
	}
	return tutors, errMsg
}

func ListTutors(db *sql.DB) []Tutor {

	results, err := db.Query("SELECT TutorID, Name, Description FROM Tutor")

	if err != nil {
		panic(err.Error())
	}
	var tutors []Tutor
	for results.Next() {
		var getTutor Tutor
		err = results.Scan(&getTutor.TutorID, &getTutor.Name, &getTutor.Description)
		if err != nil {

			panic(err.Error())
		}
		tutors = append(tutors, getTutor) //Store them in a list and use if required. --> var trips []Trip
	}
	return tutors
}

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
				if newTutor.TutorID == "" || newTutor.Name == "" || newTutor.Description == "" || newTutor.Password == "" {
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
				if tutors.Name == "" || tutors.Description == "" || tutors.Password == "" {
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
	if r.Method == "GET" { //its working
		TutorID := r.URL.Query().Get("TutorID")
		fmt.Println("TutorID: ", TutorID)
		tutors := ListTutors(db)

		json.NewEncoder(w).Encode(&tutors)
	}
	//---Deny any deletion of tutor's account or other tutor's information
	if r.Method == "DELETE" {
		TutorID := r.URL.Query().Get("TutorID")
		fmt.Println("TutorID: ", TutorID)
		tutors, errMsg := DeleteTutor(db, TutorID)

		if errMsg == "Tutor ID not found." {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Tutor's account not found"))
		} else {
			fmt.Println(tutors)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Tutor's account deleted"))
		}
	}
}

//Function main for testing purposes only
func main() {
	// instantiate tutor
	tutors, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:8012)/tutors")
	db = tutors
	// handle error
	if err != nil {
		panic(err.Error())
	}
	//handle the API connection across tutor's microservices
	router := mux.NewRouter()
	router.HandleFunc("/tutors", tutor).Methods("GET", "POST", "PUT", "DELETE")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	fmt.Println("Tutors microservice API --> Listening at port 8011")
	log.Fatal(http.ListenAndServe(":8011", handlers.CORS(originsOk, headersOk, methodsOk)(router)))

	defer db.Close()
}
