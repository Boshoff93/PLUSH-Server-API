package main

import (
  "fmt"
  "github.com/gocql/gocql"
  "net/http"
  "encoding/json"
)



func addPost(w http.ResponseWriter, r *http.Request){
  var post Post

  session := getSession()
  defer session.Close()

  if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  finished := make(chan bool)
  go func() {
    if err := session.Query("INSERT INTO posts (user_id, post_id, content) VALUES (?,?,?)",post.User_Id, post.Post_Id , post.Post).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    //Convert string uuidv1 to uuidv1 then extract the Time before sending it back to the client
    tempUUID, err := gocql.ParseUUID(post.Post_Id);
    if err != nil {
		   fmt.Printf("Something went wrong: %s", err.Error())
	  }
    var postAdded PostAdded
    postAdded.User_Id = post.User_Id
    postAdded.Post_Time = tempUUID.Time()
    postAdded.Post_Id = post.Post_Id
    postAdded.Post = post.Post
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(postAdded)
    finished <- true
  }()
  <- finished
}
