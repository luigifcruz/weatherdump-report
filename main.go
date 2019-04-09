package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
	"net/http"
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
	PlatformArch     string `schema:"platformArch,required"`
	PlatformOs       string `schema:"platformOs,required"`
	SystemCpusNumber string `schema:"systemCpusNumber,required"`
	SystemCpus       string `schema:"systemCpus,required"`
	SystemMemory     string `schema:"systemMemory,required"`
	Timestamp        string `schema:"timestamp,required"`
	VersionGUI       string `schema:"versionGUI,required"`
	VersionCLI       string `schema:"versionCLI,required"`
	Crash            string `schema:"crash,required"`
}

func crashReporter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	var c Crash
	if err := decoder.Decode(&c, r.PostForm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error %s", err)
		return
	}

	buf, _ := json.Marshal(c)

	crashID := uuid.Must(uuid.NewRandom())
	err := ioutil.WriteFile("./crashes/"+crashID.String()+"."+c.VersionCLI+".json", []byte(buf), 0755)
	err = ioutil.WriteFile("./crashes/"+crashID.String()+"."+c.VersionCLI+".txt", []byte(c.Crash), 0755)
	if err != nil {
		log.Printf("Unable to write file: %v", err)
	}

	log.Println("New crash report!", crashID.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}
