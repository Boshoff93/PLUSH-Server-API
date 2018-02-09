package main

import (
  "net/http"
  "github.com/gocql/gocql"
  "fmt"
)

func main() {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "plush_keyspace"
    session, err := cluster.CreateSession()
    if err != nil {
      fmt.Println(err);
    }
    defer session.Close()

    router := NewRouter(session)

    router.Handle("user add", addUser)
    router.Handle("user find", findUser)
    router.Handle("post add", addPost)
    router.Handle("posts get", getPosts)
    router.Handle("post delete", deletePost)
    router.Handle("user get", getUserView)
    router.Handle("profile picture add", addProfilePicture)
    router.Handle("profile picture get", getProfilePicture)
    http.Handle("/", router)
    http.ListenAndServe(":4000", nil)
}
