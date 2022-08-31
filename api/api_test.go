package main

import (
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

func setup() map[string]interface{} {

	env := make(map[string]interface{})
	rootTestDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory. %v\n", err)
	}

	if strings.Contains(strings.ToLower(runtime.GOOS), "window") {
		rootTestDir = rootTestDir + "\\..\\mock\\"
	} else {
		rootTestDir = rootTestDir + "/../mock/"
	}
	env["rootDirectory"] = rootTestDir
	env["directoriesToWatch"] = []string{"A", "B", "C"}
	env["numOfDirs"] = 4
	env["numOfFiles"] = 4
	return env
}

func TestGetDirectoriesToWatch(t *testing.T) {
	testRootDirectory := "./config/testDirectoriesToWatch.json"
	got, err := getDirectoriesToWatch(testRootDirectory)
	want := []string{"A", "B", "C"}

	if err != nil {
		t.Errorf("ERROR: Failed to get directories.")
		return
	}

	for i, v := range got {
		g := got[i]
		w := want[i]
		if strings.Compare(g, w) != 0 {
			t.Errorf("got[%v] is %v, but want[%v] is %v.", i, v, i, v)
		}
	}
}

func TestGetAllFiles(t *testing.T) {
	env := setup()
	fm := new(FileManagement)

	fm.rootDirectory = env["rootDirectory"].(string)
	fm.directoriesToWatch = env["directoriesToWatch"].([]string)

	got, err := getAllFiles(fm)
	if err != nil {
		t.Errorf("getAllFiles returned an error: %v\n", err)
	}

	numOfDirs := 0
	numOfFiles := 0
	var loopThroughMaps func(got map[string]interface{})
	loopThroughMaps = func(got map[string]interface{}) {
		for k, v := range got {
			numOfDirs += 1
			switch v.(type) {
			case []string:
				for _, i := range v.([]string) {
					numOfFiles += 1
					if !strings.HasSuffix(i, ".txt") {
						t.Errorf("Value should end with .txt, but got %v\n", i)
					}
				}
			case map[string]interface{}:
				loopThroughMaps(v.(map[string]interface{}))
			default:
				if !strings.HasSuffix(k, "C\\") {
					t.Errorf("Expected folder C to be here, but got %v with type %T\n", k, v)
				}
			}
		}
	}
	loopThroughMaps(got)

	if numOfDirs != env["numOfDirs"] {
		t.Errorf("Expected %v directories, but got %v\n", env["numOfDirs"], numOfDirs)
	}
	if numOfFiles != env["numOfFiles"] {
		t.Errorf("Expected %v directories, but got %v\n", env["numOfFiles"], numOfFiles)
	}
}
