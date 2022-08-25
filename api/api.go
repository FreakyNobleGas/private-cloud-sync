package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// FilesJSON Holds all file information to be converted to JSON
type FilesJSON struct {
	Files map[string]interface{}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func (fm *FileManagement) httpGetAllFiles(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
	files, err := getAllFiles(fm)
	if err != nil {
		log.Fatal(err)
	}

	allFilesJSON, err := json.Marshal(&FilesJSON{
		files,
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Write(allFilesJSON)
}

func main() {
	fm, err := NewFileManagement("", "", []string{})
	if err != nil {
		log.Fatalf("ERROR: Could not create NewFileManagement object. %v\n", err)
	}

	log.Println("INFO: Serving HTTP Service")
	http.HandleFunc("/getAllFiles", fm.httpGetAllFiles)
	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Print(err)
	}
}
