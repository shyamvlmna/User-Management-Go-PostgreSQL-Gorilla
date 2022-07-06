package router

import (
	"crudapp/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	//user accessible routes
	router.HandleFunc("/", controllers.UserIndex)

	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	router.HandleFunc("/user", controllers.UserHome)
	router.HandleFunc("/update/{id}", controllers.UserEdit).Methods("GET") //get the edit page
	router.HandleFunc("/update", controllers.UserUpdate).Methods("POST")   //update the edits
	router.HandleFunc("/signup", controllers.SignUp).Methods("GET")        //get the signup page
	router.HandleFunc("/addUser", controllers.UserSignUp).Methods("POST")  //add the user to db from signup page
	router.HandleFunc("/logout", controllers.UserLogout)

	//admin page path
	router.HandleFunc("/admin", controllers.AdminIndex).Methods("GET")

	//router for admin
	adminRouter := router.PathPrefix("/admin").Subrouter()

	//routes only accessible by the admin
	adminRouter.HandleFunc("/login", controllers.AdminLogin).Methods("POST")
	adminRouter.HandleFunc("/adminHome", controllers.AdminHome)
	adminRouter.HandleFunc("/view", controllers.ViewUsers).Methods("GET")    //list all users
	adminRouter.HandleFunc("/create", controllers.CreateUser).Methods("GET") //get the user create user page
	adminRouter.HandleFunc("/insertUser", controllers.InsertUser).Methods("POST")         //add the user to db from create user page
	adminRouter.HandleFunc("/delete/{id}", controllers.DeleteUser) //delete user by id
	adminRouter.HandleFunc("/logout", controllers.AdminLogout)

	return router
}
