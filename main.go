// this contains the main function which starts the server and initializes the mongo connection
package main

//importing required packages
import (
	"context"
	"log"
	"net/http"
	"os"
	"rest-api/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// global variable to store mongo client
var mongoClient *mongo.Client

// init is automatically executed function when the program starts
func init() {
	//load .env file
	//'godotenv' package is used to load the env file
	err := godotenv.Load()

	//if error occurs while loading env file, then log the error and exit the program
	if err != nil {
		log.Fatalf("error while loading env variables: %v", err)
	}

	//log the success message
	log.Println("loaded env file")

	//create mongo client
	//mongo.Connect function is used to connect to the mongo server
	//os.Getenv("MONGO_URI") is used to get the mongo uri from the env file
	//options.Client().ApplyURI is used to set the uri to the client
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatalf("error while connecting to mongo: %v", err)
	}

	//ping the mongo server using mongoClient.Ping function
	//readpref.Primary() is used to set the read preference to primary
	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("error while pinging to mongo: %v", err)
	}

	//log the success message
	log.Println("connected to mongo")
}

func main() {
	//disconnect mongo connection
	//ensures that the mongo connection is closed when the program exits
	defer mongoClient.Disconnect(context.Background())

	//get the database and collection from the mongo client
	coll := mongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	//create employee usecase
	//ensure that it can interact with the database
	empUsecase := usecase.EmployeeUsecase{MongoCollection: coll}

	//create a new router to handle the incoming requests
	r := mux.NewRouter()

	//to check if the server is running
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	//create routes for the employee
	//calls the InsertEmployee function when the request is POST
	r.HandleFunc("/employee", empUsecase.InsertEmployee).Methods(http.MethodPost)
	//calls the GetEmployeeById function when the request is GET
	r.HandleFunc("/employee/{eid}", empUsecase.GetEmployeeById).Methods(http.MethodGet)
	//calls the FindAllEmployees function when the request is GET
	r.HandleFunc("/employee", empUsecase.FindAllEmployees).Methods(http.MethodGet)
	//calls the UpdateEmployeeById function when the request is PUT
	r.HandleFunc("/employee/{eid}", empUsecase.UpdateEmployeeById).Methods(http.MethodPut)
	//calls the DeleteEmployeeById function when the request is DELETE
	r.HandleFunc("/employee/{eid}", empUsecase.DeleteEmployeeById).Methods(http.MethodDelete)
	//calls the DeleteAllEmployees function when the request is DELETE
	r.HandleFunc("/employee", empUsecase.DeleteAllEmployees).Methods(http.MethodDelete)

	//uses router 'r' to handle the incoming requests
	log.Println("server started at :3000")
	http.ListenAndServe(":3000", r)
}

// initially to check if the server is running
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running..."))
}
