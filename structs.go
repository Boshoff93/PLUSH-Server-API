package main

import (
    "time"
)

type User struct {
  User_Id     string            `json: "user_id"    cql:"uuid"`
  Firstname   string            `json: "firstname"  cql:"text"`
  Lastname    string            `json: "lastname"   cql:"text"`
  Email       string            `json: "email"      cql:"text"`
  Created_At  string            `json: "created_at" cql:"timeuuid"`
}

type SearchUser struct {
  Search        string          `cql:"text"`
}

type SearchedUsers struct {
  User_Ids     []string         `json: "user_ids"`
  Emails       []string         `json: "emails"`
  Fullnames    []string         `json: "fullnames"`
}

type Post struct {
  User_Id     string            `json: "user_id" cql:"uuid"`
  Post_Id     string            `json: "post_id" cql:"timeuuid"`
  Post        string            `json: "user_id" cql:"text"`
}

type PostAdded struct {
  User_Id       string          `json: "user_id"`
  Post_Time     time.Time       `json: "post_time"`
  Post_Id       string          `json: "post_id"`
  Post          string          `json: "post"`
}

type Posts struct {
  Post_Ids     []string         `json: "post_ids"`
  Post_Times   []time.Time      `json: "post_times"`
  Posts        []string         `json: "posts"`
}

type Blob struct {
	User_Id  string              `cql:"uuid"`
	Data     string              `cql:"text"`// Formatted as Base64 but I would prefer Base64Url...
}

type GoogleUser struct {
    Sub string `json:"sub"`
    Name string `json:"name"`
    GivenName string `json:"given_name"`
    FamilyName string `json:"family_name"`
    Profile string `json:"profile"`
    Picture string `json:"picture"`
    Email string `json:"email"`
    EmailVerified string `json:"email_verified"`
    Gender string `json:"gender"`
}

type Message struct {
  Name string                   `json:"name"`
  Data interface{}              `json:"data"`
}
