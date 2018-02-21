package main

import (
  "fmt"
  "github.com/gocql/gocql"
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
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

    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte("MyFancySecret"), nil
    },
    // When set, the middleware verifies that tokens are signed with the specific signing algorithm
    // If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
    // Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
    SigningMethod: jwt.SigningMethodHS256,
  })

    router.HandleFunc("/plush-api/login",login).Methods("POST")
    // router.HandleFunc("/plush-api/user", addUser).Methods("POST")
    // router.HandleFunc("/plush-api/user/{email}", findUser).Methods("GET")
    router.HandleFunc("/plush-api/userViewId/{user_id}", getUserViewByUserId).Methods("GET")
    //router.HandleFunc("/plush-api/userViewEmail/{email}", getUserViewByEmail).Methods("GET")
    router.HandleFunc("/plush-api/searchUsers/{like_name}", searchUsers).Methods("GET")
    router.HandleFunc("/plush-api/post", addPost).Methods("POST")
    router.HandleFunc("/plush-api/post", deletePost).Methods("DELETE")
    router.HandleFunc("/plush-api/getposts/{user_id}", getPosts).Methods("GET")
    router.HandleFunc("/plush-api/profilePicture", addProfilePicture).Methods("POST")
    router.HandleFunc("/plush-api/profilePicture/{user_id}", getProfilePicture).Methods("GET")
    http.ListenAndServe(":8000", handlers.CORS(headersOk, methodsOk, originsOk)(loggedRouter))

}
