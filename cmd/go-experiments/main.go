package main

import (
  "log"
  "fmt"
  "net/http"
  "os"

  "github.com/gorilla/mux"
  "github.com/DimitarYankov/go-experiments/internal/diagnostics"
)

func main() {
   log.Print("Starting the application... ")

   blPort := os.Getenv("PORT")
   if (len(blPort) == 0) {
     log.Fatal("The app port should be set")
   }
   diagnosticPort := os.Getenv("DIAG_PORT")
   if (len(diagnosticPort) == 0) {
     log.Fatal("The daignostic port should be set")
   }
   router := mux.NewRouter()
   router.HandleFunc("/", hello)

   go func () {
     err := http.ListenAndServe(":"+blPort, router)
     if err != nil {
       log.Fatal(err)
     }
   }()

   diagnostics := diagnostics.NewDiagnostics()
   err := http.ListenAndServe(":"+diagnosticPort, diagnostics)
   if err != nil {
     log.Fatal(err)
   }
}

func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, http.StatusText(http.StatusOK))
}
