package main

import (
	"os"

	"github.com/codegangsta/negroni"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
    "encoding/json"
    "log"
    "github.com/tobyjsullivan/life/signup-svc/profile"
)

var logger *log.Logger
var profileSvc *profile.Service

func init() {
    logger = log.New(os.Stdout, "[signup-svc] ", 0)

    profileSvc = profile.NewService()

    logger.Println("Initialized.")
}

func main() {
	r := buildRoutes()

	n := negroni.New()
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n.Run(":" + port)
}

func buildRoutes() http.Handler {
	r := mux.NewRouter()
    r.HandleFunc("/", statusHandler).
            Methods("GET")
    r.HandleFunc("/commands/create-profile", createProfileHandler).
            Methods("POST")
    r.HandleFunc("/commands/change-name", changeNameHandler).
            Methods("POST")

	return r
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Online.")
}

func createProfileHandler(w http.ResponseWriter, r *http.Request) {
    logger.Println("Command received: Create Profile")
    decoder := json.NewDecoder(r.Body)
    var req createProfileReq
    err := decoder.Decode(&req)
    if err != nil {
        logger.Println("Failed to parse request: "+err.Error())
        http.Error(w, "Failed to parse request: "+err.Error(), http.StatusBadRequest)
        return
    }

    logger.Println("Calling createProfile")
    go profileSvc.CreateProfile(req.FirstName, req.LastName)

    logger.Println("Respoding to request")

    w.WriteHeader(http.StatusAccepted)
    fmt.Fprint(w, "")
    logger.Println("Command Accepted.")
}

type createProfileReq struct {
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
}

func changeNameHandler(w http.ResponseWriter, r *http.Request) {
    logger.Println("Command received: Change Name")
    decoder := json.NewDecoder(r.Body)
    var req changeNameReq
    err := decoder.Decode(&req)
    if err != nil {
        logger.Println("Failed to parse request: "+err.Error())
        http.Error(w, "Failed to parse request: "+err.Error(), http.StatusBadRequest)
        return
    }

    logger.Println("Calling createProfile")
    go profileSvc.ChangeName(req.ProfileID, req.FirstName, req.LastName)

    logger.Println("Respoding to request")

    w.WriteHeader(http.StatusAccepted)
    fmt.Fprint(w, "")
    logger.Println("Command Accepted.")
}

type changeNameReq struct {
    ProfileID string `json:"profileId"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
}