package main

import (
  "encoding/json"
  "fmt"
  "github.com/gocql/gocql"
  "net/http"
)



func addPost(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var post Post
  if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  finished := make(chan bool)
  if(post.Type_Of_Post == 1) {
    post.Post = post.Post_Id + "_post_picture"
  }
  go func() {



    if err := session.Query("INSERT INTO posts (user_id, post_id, content, type) VALUES (?,?,?,?)",post.User_Id, post.Post_Id , post.Post, post.Type_Of_Post).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      return
    }

    if err := session.Query("INSERT INTO posts_likes_dislikes (post_id, user_id, like, dislike) VALUES (?,?,?,?)",
                                          post.Post_Id, post.User_Id , 0, 0).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
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
    postAdded.Type_Of_Post = post.Type_Of_Post
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(postAdded)
    finished <- true
  }()
  <- finished
}
