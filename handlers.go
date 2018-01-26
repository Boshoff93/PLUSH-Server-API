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
    client.send <- Message{"error", "could not decode"}
    return
  }

  go func() {
    _, err := r.Table("user").
      Insert(user).
      RunWrite(client.session)
    if err != nil {
      fmt.Println("failed")
      client.send <- Message{"error", err.Error()}
    }
    client.send <- Message{"user add", user}
  }()

}
