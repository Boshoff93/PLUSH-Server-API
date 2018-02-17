package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
  "github.com/gocql/gocql"
)

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
    
    router.Handle("post add", addPost)
    http.Handle("/", router)
    http.ListenAndServe(":4000", nil)
}

func addPost(client *Client, data interface{}){
  var post Post
  err := mapstructure.Decode(data, &post)
  if err != nil {
    client.send <- Message{"error", "could not decode addPost"}
    return
  }
  go func() {
    if err := client.session.Query("INSERT INTO posts (user_id, post_id, content) VALUES (?,?,?)",post.User_Id, post.Post_Id , post.Post).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    //Convert string uuidv1 to uuidv1 then extract the Time before sending it back to the client
    tempUUID, err := gocql.ParseUUID(post.Post_Id);
    if err != nil {
       client.send <- Message{"error", err.Error()}
		   fmt.Printf("Something went wrong: %s", err)
	  }
    var postAdded PostAdded
    postAdded.User_Id = post.User_Id
    postAdded.Post_Time = tempUUID.Time()
    postAdded.Post_Id = post.Post_Id
    postAdded.Post = post.Post
    client.send <- Message{"post add", postAdded}
  }()

}
