package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FilesJSON Holds all file information to be converted to JSON
type FilesJSON struct {
	Files map[string]interface{}
}

func getFiles(nextDir string, allFiles map[string]interface{}) (map[string]interface{}, error) {
	if allFiles[nextDir] == nil {
		allFiles[nextDir] = make([]string, 0)

	}
	var nextDirMapping = allFiles[nextDir].([]string)

	err := filepath.Walk(nextDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			if strings.Compare(nextDir, path) != 0 {
				allFiles[path] = make(map[string]interface{})
				getFiles(path, allFiles[path].(map[string]interface{}))
			}
		} else {
			nextDirMapping = append(nextDirMapping, path)
		}

		return nil
	})

	allFiles[nextDir] = nextDirMapping
	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}

	return allFiles, err
}

func getAllFiles() (map[string]interface{}, error) {
	var allFiles = make(map[string]interface{})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("ERROR: Can't get home directory path.")
	}
	// TODO: Remove Test Dir
	homeDir = "C:\\Users\\nickq\\Desktop\\testing"

	getFolders := []string{"Documents", "Pictures", "Videos"}
	getFolders = []string{"A", "B", "C"}
	for _, f := range getFolders {
		nextDir := homeDir + "\\" + f
		allFiles[nextDir] = make([]string, 0)
		allFiles, err = getFiles(nextDir, allFiles)
	}

	var printer func(allFiles map[string]interface{})
	printer = func(allFiles map[string]interface{}) {
		for k, v := range allFiles {
			switch v := v.(type) {
			case map[string]interface{}:
				printer(v)
			case []string:
				fmt.Printf("Retrieved %v files in %v.\n", len(v), k)
			default:
				fmt.Printf("I don't know about type %T!\n", v)
			}
		}
	}
	printer(allFiles)

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

	allFilesJSON, err := json.Marshal(&FilesJSON{
		files,
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Write(allFilesJSON)
}

func main() {
	log.Println("INFO: Testing getAllFiles()")
	_, _ = getAllFiles()
	log.Println("INFO: Serving HTTP Service")
	http.HandleFunc("/getAllFiles", httpGetAllFiles)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Print(err)
	}
}
