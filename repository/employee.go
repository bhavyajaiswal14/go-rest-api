// responsible for interacting with the database to perform CRUD operations
package repository

//import the necessary packages
import (
	"context"
	"fmt"
	"log"
	"rest-api/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository struct {
	//MongoCollection is a pointer to the mongo.Collection object
	//that represents the collection in the MongoDB database
	MongoCollection *mongo.Collection
}

// adds a new employee to the database
// InsertEmployee method takes an employee object as input and returns the inserted ID or an error
func (er *EmployeeRepository) InsertEmployee(employee *schema.Employee) (interface{}, error) {

	//'result' contains ID of the inserted document
	//'err' contains any error that occurred during the insertion
	//'InsertOne' method inserts a single document into the collection
	//'context.Background()' means just run the operation without any timeout
	result, err := er.MongoCollection.InsertOne(context.Background(), employee)

	//if an error occurred during the insertion, log the error and return it
	if err != nil {
		log.Printf("error while inserting employee: %v", err)
		return nil, err
	}

	//return the ID of the inserted document
	return result.InsertedID, nil
}

// GetEmployeeById method takes an employee ID as input and returns the employee object or an error
// 'eid string' is given as input to get the employee by ID
// output is employee object or an error
func (er *EmployeeRepository) GetEmployeeById(eid string) (*schema.Employee, error) {
	var employee schema.Employee //to store the employee object

	//FindOne method returns a single document that matches the filter
	//Decode method decodes the document into the employee object
	err := er.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "e_id", Value: eid}}).Decode(&employee)

	//if an error occurred while getting the employee, log the error and return it
	if err != nil {
		log.Printf("error while getting employee: %v", err)
		return nil, err
	}

	//return the employee object
	return &employee, nil
}

// FindAllEmployees method returns all the employees in the database
// output is a slice of employee objects or an error
func (er *EmployeeRepository) FindAllEmployees() ([]schema.Employee, error) {
	//Find method returns all the documents in the collection
	results, err := er.MongoCollection.Find(context.Background(), bson.D{})

	//if an error occurred while getting the employees, log the error and return it
	if err != nil {
		log.Printf("error while getting employees: %v", err)
		return nil, err
	}

	//decode the documents into a list or slice of employee objects
	var employees []schema.Employee
	//All method decodes all the documents into the slice of employee objects
	err = results.All(context.Background(), &employees)

	//if an error occurred while decoding the employees, log the error and return it
	if err != nil {
		return nil, fmt.Errorf("error while decoding employees: %s", err.Error())
	}

	//return the list of employee objects
	return employees, nil
}

// UpdateEmployeeById method takes an employee ID and an updated employee object as input
// and returns the number of modified documents or an error
// 'eid string' is given as input to update the employee by ID
func (er *EmployeeRepository) UpdateEmployeeById(eid string, updateEmp *schema.Employee) (int64, error) {
	//UpdateOne method updates a single document that matches the filter
	result, err := er.MongoCollection.UpdateOne(context.Background(),

		//filter to find the document to update
		bson.D{{Key: "e_id", Value: eid}},
		//'$set' updates the fields specified and leaves the other fields unchanged
		bson.D{{Key: "$set", Value: updateEmp}})

	//if an error occurred while updating the employee, log the error and return it
	if err != nil {
		log.Printf("error while updating employee: %v", err)
		return 0, err
	}

	//return the number of modified documents
	return result.ModifiedCount, nil
}

// DeleteEmployeeById method takes an employee ID as input and returns the number of deleted documents or an error
func (er *EmployeeRepository) DeleteEmployeeById(eid string) (int64, error) {

	//DeleteOne method deletes a single document that matches the filter
	result, err := er.MongoCollection.DeleteOne(context.Background(),

		//filter to find the document to delete
		bson.D{{Key: "e_id", Value: eid}})

	//if an error occurred while deleting the employee, log the error and return it
	if err != nil {
		log.Printf("error while deleting employee: %v", err)
		return 0, err
	}

	//return the number of deleted documents
	return result.DeletedCount, nil
}

// DeleteAllEmployees method deletes all the employees in the database
func (er *EmployeeRepository) DeleteAllEmployees() (int64, error) {

	//DeleteMany method deletes all the documents in the collection
	result, err := er.MongoCollection.DeleteMany(context.Background(), bson.D{})

	//if an error occurred while deleting all the employees, log the error and return it
	if err != nil {
		log.Printf("error while deleting all employees: %v", err)
		return 0, err
	}

	//return the number of deleted documents
	return result.DeletedCount, nil
}
