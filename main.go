package main

import (
  "net/http"
  r "gopkg.in/gorethink/gorethink.v4"
  "log"
)

func main() {
    session, err := r.Connect(r.ConnectOpts{
      Address: "localhost:28015",
      Database: "rtsupport",
    })
    if err != nil {
      log.Panic(err.Error())
    }

    router := NewRouter(session)

    router.Handle("user add", addUser)
    router.Handle("user find", findUser)
    http.Handle("/", router)
    http.ListenAndServe(":4000", nil)
}
