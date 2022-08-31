package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileManagement struct {
	directoriesToWatchConfig string
	directoriesToWatch       []string
	filesByDirectory         map[string]interface{}
	filesMetaData            map[string]map[string]string
	rootDirectory            string
}

func NewFileManagement(
	rootDirectory string,
	directoriesToWatchConfig string,
	directoriesToWatch []string,
) (*FileManagement, error) {

	log.Println("INFO: Constructing FileManagement object.")
	var fm = new(FileManagement)
	var err error

	if rootDirectory != "" {
		fm.rootDirectory = rootDirectory
	} else {
		fm.rootDirectory, err = os.UserHomeDir()
		if err != nil {
			return nil, err
		}
	}

	if directoriesToWatchConfig != "" {
		fm.directoriesToWatchConfig = "./config/directoriesToWatch.json"
	} else {
		fm.directoriesToWatchConfig = directoriesToWatchConfig
	}

	if len(directoriesToWatch) != 0 {
		fm.directoriesToWatch = directoriesToWatch
	} else {
		fm.directoriesToWatch, err = getDirectoriesToWatch(fm.directoriesToWatchConfig)
		if err != nil {
			return nil, err
		}
	}

	fm.filesByDirectory, err = getAllFiles(fm)
	if err != nil {
		return nil, err
	}

	// TODO Implement constructor for filesMetaData

	return fm, nil
}

func getDirectoriesToWatch(fileName string) ([]string, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("ERROR: Unable to retreive directories to watch configuration.")
		return nil, err
	}

	var configMap map[string][]string
	err = json.Unmarshal(f, &configMap)
	if err != nil {
		log.Printf("ERROR: Unable to parse JSON config file.\n")
		return nil, err
	}

	directories, exist := configMap["directories"]
	if exist {
		if len(directories) > 0 {
			return directories, nil
		}
		return nil, fmt.Errorf("ERROR: No directories are listed in the directoriesToWatch.json config")
	}
	return nil, fmt.Errorf("ERROR: \"directories\" key does not exist in directoriesToWatch.json config")
}

func getFiles(nextDir string, allFiles map[string]interface{}) (map[string]interface{}, error) {

	var currentDirectory = nextDir
	err := filepath.Walk(nextDir, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if allFiles[nextDir] == nil {
			allFiles[nextDir] = make([]string, 0)
		}

		if info.IsDir() {
			if strings.Compare(nextDir, path) != 0 {
				if allFiles[path] == nil {
					allFiles[path] = make([]string, 0)
				}
				currentDirectory = path
			}
		} else {
			allFiles[currentDirectory] = append(allFiles[currentDirectory].([]string), path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}

	return allFiles, err
}

func getAllFiles(fm *FileManagement) (map[string]interface{}, error) {
	var allFiles = make(map[string]interface{})
	var err error

	for _, f := range fm.directoriesToWatch {
		nextDir := fm.rootDirectory + "\\" + f
		allFiles[nextDir] = make([]string, 0)
		allFiles, err = getFiles(nextDir, allFiles)
	}

	// TODO Move to own function
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
