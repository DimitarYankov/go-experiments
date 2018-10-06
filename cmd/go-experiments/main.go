package main

import (
  "log"
  "fmt"
  "net/http"
  "os"

  "github.com/gorilla/mux"
  "github.com/DimitarYankov/go-experiments/internal/diagnostics"
)

type serverConf struct{
  port string
  router http.Handler
  name string
}

func main() {
   log.Print("Starting the application... ")

   blPort := os.Getenv("PORT")
   if (len(blPort) == 0) {
     log.Fatal("The app port should be set")
   }
   router := mux.NewRouter()
   router.HandleFunc("/", hello)

   diagnosticPort := os.Getenv("DIAG_PORT")
   if (len(diagnosticPort) == 0) {
     log.Fatal("The daignostic port should be set")
   }
   diagnostics := diagnostics.NewDiagnostics()

   possibleErrors := make(chan error, 2)
   configurations := []serverConf{
     {
       port: blPort,
       router: router,
       name: "application server",
     },
     {
       port: diagnosticPort,
       router: diagnostics,
       name: "diagnostic server",
     },
   }

   servers := make([]*http.Server, 2)

   for i, c := range configurations {
     go func (conf serverConf, i int) {
       log.Printf("The %s is preparing to handle connections...", conf.name)
       servers[i] = &http.Server{
         Addr: ":"+conf.port,
         Handler: conf.router,
       }
       err := servers[i].ListenAndServe()
       if (err != nil) {
         possibleErrors <- err
       }
     }(c, i)
   }

   select {
   case err := <-possibleErrors:
     // for _, s := range servers {
     //   // context.Background()
     //   // s.Shutdown()
     // }
     log.Fatal(err)

   }

}

func hello(w http.ResponseWriter, r *http.Request) {
  log.Print("The hello handler was called")
  fmt.Fprint(w, http.StatusText(http.StatusOK))
}
