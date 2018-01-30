package main

import (
  "github.com/mitchellh/mapstructure"
  r "gopkg.in/gorethink/gorethink.v4"
  "fmt"
)

type User struct {
  Id    string `gorethink:"id"`
  Name  string  `gorethink:"name"`
}

type Post struct {
  User_Id    string `gorethink:"user_id"`
  Post  string  `gorethink:"post"`
}

type Posts struct {
  User_Id    string `gorethink:"user_id"`
  Posts  []string `gorethink:"posts"`
}

type Message struct {
  Name string `json:"name"`
  Data interface{} `json:"data"`
}

func findUser(client *Client, data interface{}){
  var user User
  fmt.Println(data)
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", err.Error()}
    return
  }

  go func() {
    res, err := r.Table("user").Run(client.session)
    if err != nil {
      client.send <- Message{"error", "can't access table"}
      return
    }
    var row User
    for res.Next(&row) {
      if(row.Name == user.Name) {
        client.send <- Message{"username unavailible", row}
        res.Close()
        return
      }
    }
    fmt.Println(user.Name + " is not taken")
    client.send <- Message{"username availible", user}
  }()

}

func addUser(client *Client, data interface{}){
  var user User
  fmt.Println(data)
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode user"}
    return
  }

  go func() {
    _, err := r.Table("user").
      Insert(user).
      RunWrite(client.session)
    if err != nil {
      fmt.Println("failed to add initial user")
      client.send <- Message{"error", err.Error()}
    }

    var posts Posts
    posts.User_Id = user.Id
    posts.Posts = []string{}

    _, err2 := r.Table("post").
      Insert(posts).
      RunWrite(client.session)
    if err2 != nil {
      fmt.Println("failed to add initial post")
      client.send <- Message{"error", err.Error()}
    }

    client.send <- Message{"user add", user}
  }()

}

func addPost(client *Client, data interface{}){
  var post Post
  fmt.Println(data)
  err := mapstructure.Decode(data, &post)
  if err != nil {
    client.send <- Message{"error", "could not decode post"}
    return
  }

  go func() {
    _, err := r.Table("post").Filter(r.Row.Field("user_id").Eq(post.User_Id)).Update(map[string]interface{}{"posts": r.Row.Field("posts").Append(post.Post)}).
      RunWrite(client.session)
    if err != nil {
      fmt.Println("failed to update posts")
      client.send <- Message{"error", err.Error()}
    }
    client.send <- Message{"post add", post}
  }()

}

func getPosts(client *Client, data interface{}){
  var user User
  fmt.Println(data)
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode"}
    return
  }

  go func() {
    res, err := r.Table("post").Run(client.session)
    if err != nil {
      client.send <- Message{"error", "can't access table"}
      return
    }
    var row Posts
    for res.Next(&row) {
      if(row.User_Id == user.Id) {
        res.Close()
        fmt.Println("------------")
        fmt.Println(row)

        client.send <- Message{"posts get", row.Posts}
        return
      }
    }

  }()

}
