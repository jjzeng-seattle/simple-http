package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
  build:= os.Getenv("BUILD")
  if build == "" {
    build = "unknown build"
  }

  echo, ok := r.URL.Query()["echo"]
  if ok {
	  fmt.Fprintf(w, "From %s: You sent %s!", build, echo[0])
  } else {
	  fmt.Fprintf(w, "From %s: Hello!", build)
  }
}

func healthcheck_handler(w http.ResponseWriter, r *http.Request) {
  status, ok := r.URL.Query()["status"]
  // If the request doesn't contain status, treat it as success.
  if !ok || status[0] != "f" {
    fmt.Fprint(w, "OK")
  } else {
     w.WriteHeader(http.StatusInternalServerError)
     w.Write([]byte("500 - Something bad happened!"))
  }
}

func main() {
  log.Print("helloworld: starting server...")

  // The caller can control the status of response by setting params.
  // "status=f" will return 500, otherwise return 200.
  // This is for an easy way to simulate failure.
  http.HandleFunc("/healthcheck", healthcheck_handler)

  // If you send "echo=foo", the response will contain foo. Otherwise it will
  // just say "Hello".
  http.HandleFunc("/", handler)

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  log.Printf("helloworld: listening on port %s", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
