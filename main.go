package main

import (
  "net/http"
  "github.com/gocql/gocql"
  "github.com/gorilla/mux"
  "fmt"
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

    router.HandleFunc("/plush-api/user", addUser).Methods("POST")
    router.HandleFunc("/plush-api/user/{email}", findUser).Methods("GET")
    router.HandleFunc("/plush-api/userviewId/{user_id}", getUserViewByUserId).Methods("GET")
    router.HandleFunc("/plush-api/userviewEmail/{email}", getUserViewByEmail).Methods("GET")
    router.HandleFunc("/plush-api/searchUsers/{like_name}", searchUsers).Methods("GET")
    router.HandleFunc("/plush-api/post", addPost).Methods("POST")
    router.HandleFunc("/plush-api/post", deletePost).Methods("DELETE")
    router.HandleFunc("/plush-api/getposts/{user_id}", getPosts).Methods("GET")
    router.HandleFunc("/plush-api/profilePicture", addProfilePicture).Methods("POST")
    router.HandleFunc("/plush-api/profilePicture/{user_id}", getProfilePicture).Methods("GET")
    http.ListenAndServe(":8000", router)

}
