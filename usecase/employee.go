package usecase

//importing required packages
import (
	"encoding/json"
	"log"
	"net/http"
	"rest-api/repository"
	"rest-api/schema"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmployeeUsecase struct
// MongoCollection is a pointer to mongo.Collection
type EmployeeUsecase struct {
	MongoCollection *mongo.Collection
}

// this struct is used to send response to client
type Response struct {
	Data  interface{} `json:"data,omitempty"`  //data to be sent
	Error string      `json:"error,omitempty"` //error message if any
}

// Function to insert employee
func (svc EmployeeUsecase) InsertEmployee(w http.ResponseWriter, r *http.Request) {
	// code to insert employee
	//adds content type to specify response type
	w.Header().Add("Content-Type", "application/json")

	//stores response to be sent
	res := Response{}
	//defer is used to encode and send response to client after function completes
	defer json.NewEncoder(w).Encode(res)

	var emp schema.Employee

	//decode request body to employee struct
	err := json.NewDecoder(r.Body).Decode(&emp)
	//if error while decoding, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request body", err)
		res.Error = err.Error()
		return
	}
	//assign new employee id
	emp.EID = uuid.NewString()

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	//insert employee
	insertID, err := repo.InsertEmployee(&emp)
	//if error while inserting, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while inserting employee", err)
		res.Error = err.Error()
		return
	}
	//send response with employee id
	res.Data = emp.EID
	w.WriteHeader(http.StatusOK)

	log.Println("employee inserted with id:", insertID, emp)
}

// Function to get employee by id
func (svc EmployeeUsecase) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	// code to get employee by id
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id from request
	//'eid' extracts employee id from the request using mux.Vars(r)
	eid := mux.Vars(r)["eid"]
	log.Println("getting employee with id:", eid)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	//get employee by id
	emp, err := repo.GetEmployeeById(eid)
	//if error while getting employee, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while getting employee", err)
		res.Error = err.Error()
		return
	}
	//send response with employee details
	res.Data = emp
	w.WriteHeader(http.StatusOK)
}

// Function to get all employees
func (svc EmployeeUsecase) FindAllEmployees(w http.ResponseWriter, r *http.Request) {
	// code to get all employees
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}

	//get all employees
	emp, err := repo.FindAllEmployees()
	//if error while getting employees, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while getting employee", err)
		res.Error = err.Error()
		return
	}
	//send response with all employees
	res.Data = emp
	w.WriteHeader(http.StatusOK)
}

// Function to update employee by id
func (svc EmployeeUsecase) UpdateEmployeeById(w http.ResponseWriter, r *http.Request) {
	// code to update employee by id
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id
	eid := mux.Vars(r)["eid"]
	log.Println("updating employee with id:", eid)

	//if employee id is empty, send error response
	if eid == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid employee id")
		res.Error = "invalid employee id"
		return
	}

	var emp schema.Employee

	//decode request body to employee struct
	err := json.NewDecoder(r.Body).Decode(&emp)
	//if error while decoding, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request body", err)
		res.Error = err.Error()
		return
	}

	//assign employee id
	emp.EID = eid

	//update employee by id
	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}
	count, err := repo.UpdateEmployeeById(eid, &emp)
	//if error while updating employee, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while updating employee", err)
		res.Error = err.Error()
		return
	}
	//send response with count of updated employees
	res.Data = count
	w.WriteHeader(http.StatusOK)
}

// Function to delete employee by id
func (svc EmployeeUsecase) DeleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	// code to delete employee by id
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	//get employee id
	eid := mux.Vars(r)["eid"]
	log.Println("deleting employee with id:", eid)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}
	count, err := repo.DeleteEmployeeById(eid)
	//if error while deleting employee, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while deleting employee", err)
		res.Error = err.Error()
		return
	}
	//send response with count of deleted employees
	res.Data = count
	w.WriteHeader(http.StatusOK)
}

// Function to delete all employees
func (svc EmployeeUsecase) DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	// code to delete all employees
	w.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(w).Encode(res)

	repo := repository.EmployeeRepository{MongoCollection: svc.MongoCollection}
	//delete all employees
	count, err := repo.DeleteAllEmployees()
	//if error while deleting employees, send error response
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error while deleting employee", err)
		res.Error = err.Error()
		return
	}
	//send response with count of deleted employees
	res.Data = count
	w.WriteHeader(http.StatusOK)

}
