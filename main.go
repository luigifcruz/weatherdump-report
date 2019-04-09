package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func main() {
	color.Green("Reporter has started.")

	r := mux.NewRouter()
	r.HandleFunc("/report/crash", crashReporter)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:5415", nil))
}

type Crash struct {
	PlatformArch 	 string `schema:"platformArch,required"`
	PlatformOs 		 string `schema:"platformOs,required"`
	SystemCpusNumber int	`schema:"systemCpusNumber,required"`
    SystemCpus 		 string	`schema:"systemCpus,required"`
    SystemMemory	 int	`schema:"systemMemory,required"`
	Timestamp 		 int	`schema:"timestamp,required"`
	Crash 			 string	`schema:"crash,required"`
}

func crashReporter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	var c Crash
	if err := decoder.Decode(&c, r.PostForm); err != nil {
		
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error")
		return
	}

	buf, _ := json.Marshal(c)

	crashId := uuid.Must(uuid.NewRandom())
	err := ioutil.WriteFile("./crashes/" + crashId.String() + ".json", []byte(buf), 0755)
	err = ioutil.WriteFile("./crashes/" + crashId.String() + ".txt", []byte(c.Crash), 0755)
    if err != nil {
        log.Printf("Unable to write file: %v", err)
	}
	
	log.Println("New crash report!", crashId.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}