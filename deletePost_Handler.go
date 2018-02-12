package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
)


func deletePost(client *Client, data interface{}){
  var post Post
  err := mapstructure.Decode(data, &post)
  if err != nil {
    client.send <- Message{"error", "could not decode addPost"}
    return
  }
   go func() {
    if err := client.session.Query("DELETE FROM posts WHERE user_id = ? AND post_id = ?",post.User_Id, post.Post_Id).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    client.send <- Message{"post delete", post}
  }()
}
