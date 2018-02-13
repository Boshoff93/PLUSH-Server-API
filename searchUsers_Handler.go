package main

import (
  "github.com/mitchellh/mapstructure"
)

func searchUsers(client *Client, data interface{}){
  var search SearchUser
  err := mapstructure.Decode(data, &search)
  if err != nil {
    client.send <- Message{"error", "could not decode getUserView"}
    return
  }

  go func() {
    var searchedUsers SearchedUsers
    var user_id string
    var email string
    var fullname string
    itr := client.session.Query("SELECT email, fullname, user_id FROM users_by_email_like_fullname WHERE fullname Like ? LIMIT 50","%" + search.Search + "%").Iter()
    for itr.Scan(&email,&fullname, &user_id) {
		    searchedUsers.User_Ids = append(searchedUsers.User_Ids, user_id)
        searchedUsers.Emails = append(searchedUsers.Emails, email)
        searchedUsers.Fullnames = append(searchedUsers.Fullnames ,fullname)
      }
    client.send <- Message{"search users", searchedUsers}
  }()

}
