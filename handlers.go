package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
  "time"
)

type User struct {
  User_Id    string `cql:"uuid"`
  Name  string    `cql:"name"`
}

type Post struct {
  User_Id     string `cql:"uuid"`
  Post_Id     string `cql:"timeuuid"`
  Post        string
}

type Posts struct {
  Post_Ids     []time.Time
  Posts       []string
}

type Message struct {
  Name string `json:"name"`
  Data interface{} `json:"data"`
}

func findUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode findUser"}
    return
  }

  go func() {
    var name string
    var user_id string
    if err := client.session.Query("SELECT * FROM users WHERE name = ?",user.Name).Scan(&name, &user_id); err != nil {
      client.send <- Message{"username availible", user}
      return
    }
    user.User_Id = user_id
    fmt.Println(name + " is taken")
    client.send <- Message{"username unavailible", user}
  }()

}


func addUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode addUser"}
    return
  }

  go func() {
  if err := client.session.Query("INSERT INTO users (name, user_id) VALUES (?,?)",user.Name, user.User_Id).Exec(); err != nil {
    fmt.Println(err.Error());
  }
  client.send <- Message{"user add", user}
  }()

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
    client.send <- Message{"post add", post}
  }()

}

func getPosts(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getPosts"}
    return
  }

  go func() {
    var posts Posts
    var post_id time.Time
    var content string
    itr := client.session.Query("SELECT toTimeStamp(post_id), content FROM posts WHERE user_id = ?",user.User_Id).Iter()
    for itr.Scan(&post_id, &content) {
		    posts.Post_Ids = append(posts.Post_Ids,post_id)
        posts.Posts = append(posts.Posts ,content)
	     }
    client.send <- Message{"posts get", posts}
  }()

}
