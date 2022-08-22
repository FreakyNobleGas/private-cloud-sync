package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

type FilesJSON struct {
	Files []string
}

func getAllFiles() (allFiles []string, err error) {
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		allFiles = append(allFiles, path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return allFiles, err
	}

	return allFiles, err
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func httpGetAllFiles(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
	files, err := getAllFiles()
	if err != nil {
		log.Fatal(err)
	}

	allFiles, err := json.Marshal(&FilesJSON{
		files,
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Write(allFiles)
}

func main() {
	http.HandleFunc("/getAllFiles", httpGetAllFiles)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Print(err)
	}
}