package main

import (
  "fmt"
  "github.com/gocql/gocql"
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  "net/http"
  "os"

)

func getSession() *gocql.Session {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "plush_keyspace"
    session, err := cluster.CreateSession()
    if err != nil {
      fmt.Println(err);
    }
    return session
}

func main() {

    router := mux.NewRouter()
    loggedRouter := handlers.LoggingHandler(os.Stdout, router)
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

    router.HandleFunc("/plush-api/login",login).Methods("POST")
    // router.HandleFunc("/plush-api/user", addUser).Methods("POST")
    // router.HandleFunc("/plush-api/user/{email}", findUser).Methods("GET")
    //router.HandleFunc("/plush-api/userViewEmail/{email}", getUserViewByEmail).Methods("GET")
    router.HandleFunc("/plush-api/userViewId/{user_id}", ValidateMiddleware(getUserViewByUserId)).Methods("GET")
    router.HandleFunc("/plush-api/searchUsers/{like_name}", ValidateMiddleware(searchUsers)).Methods("GET")
    router.HandleFunc("/plush-api/post",  ValidateMiddleware(addPost)).Methods("POST")
    router.HandleFunc("/plush-api/post",  ValidateMiddleware(deletePost)).Methods("DELETE")
    router.HandleFunc("/plush-api/getposts/{user_id}",  ValidateMiddleware(getPosts)).Methods("GET")
    router.HandleFunc("/plush-api/profilePicture",  ValidateMiddleware(addProfilePicture)).Methods("POST")
    router.HandleFunc("/plush-api/profilePicture/{user_id}",  ValidateMiddleware(getProfilePicture)).Methods("GET")
    http.ListenAndServe(":8000", handlers.CORS(headersOk, methodsOk, originsOk)(loggedRouter))

}
