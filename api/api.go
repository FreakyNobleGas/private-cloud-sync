package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// FilesJSON Holds all file information to be converted to JSON
type FilesJSON struct {
	Files map[string][]string
}

func getAllFiles(allFiles *map[string][]string) (map[string][]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("ERROR: Can't get home directory path.")
	}

	getFolders := []string{"Documents", "Pictures", "Videos"}
	for _, f := range getFolders {
		nextDir := homeDir + "\\" + f
		(*allFiles)[nextDir] = make([]string, 0)
		var nextDirMapping = (*allFiles)[nextDir]
		err = filepath.Walk(nextDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			nextDirMapping = append(nextDirMapping, path)
			return nil
		})
		(*allFiles)[nextDir] = nextDirMapping
		if err != nil {
			fmt.Printf("error walking the path: %v\n", err)
			return *allFiles, err
		}
	}
	for k, v := range *allFiles {
		log.Printf("INFO: Retrieved %v files in %v.", len(v), k)
	}

	return *allFiles, err
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func httpGetAllFiles(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
	var allFiles = make(map[string][]string)
	files, err := getAllFiles(&allFiles)
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
	var allFiles = make(map[string][]string)
	t, _ := getAllFiles(&allFiles)
	for k, _ := range t {
		log.Println(k)
	}
	log.Println("INFO: Serving HTTP Service")
	http.HandleFunc("/getAllFiles", httpGetAllFiles)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Print(err)
	}
}
