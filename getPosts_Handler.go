package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
  "time"
)



func getPosts(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getPosts"}
    return
  }
  go func() {
    var posts Posts
    var post_id string
    var content string
    var post_time time.Time
    itr := client.session.Query("SELECT toTimeStamp(post_id), post_id, content FROM posts WHERE user_id = ?",user.User_Id).Iter()
    for itr.Scan(&post_time,&post_id, &content) {
		    posts.Post_Ids = append(posts.Post_Ids, post_id)
        posts.Post_Times = append(posts.Post_Times, post_time)
        posts.Posts = append(posts.Posts ,content)
      }
      fmt.Println(posts)
    client.send <- Message{"posts get", posts}
  }()

}
