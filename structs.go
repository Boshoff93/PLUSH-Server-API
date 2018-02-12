package main

import (
    "time"
)

type User struct {
  User_Id     string      `cql:"uuid"`
  Firstname   string      `cql:"firstname"`
  Lastname    string      `cql:"lastname"`
  Email       string      `cql:"email"`
  Password    string      `cql:"password"`
  Created_At  string      `cql:"timeuuid"`
}

type Post struct {
  User_Id     string `cql:"uuid"`
  Post_Id     string `cql:"timeuuid"`
  Post        string
}

type PostAdded struct {
  User_Id       string
  Post_Time     time.Time
  Post_Id       string
  Post          string
}

type Posts struct {
  Post_Ids     []string
  Post_Times   []time.Time
  Posts       []string
}

type Blob struct {
	User_Id  string    `cql:"uuid"`
	Data     string    `cql:"text"`// Formatted as Base64 but I would prefer Base64Url...
}


type Message struct {
  Name string `json:"name"`
  Data interface{} `json:"data"`
}
