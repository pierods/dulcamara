package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pierods/dulcamara"
)

const fileSuffix = ".mock"

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	if stat.Mode().IsDir() {
		return true
	}
	return false
}

func readMockFiles(path string) ([]dulcamara.Rule, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []dulcamara.Rule{}, err
	}
	var mocks []dulcamara.Rule

	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), fileSuffix) {
			continue
		}
		fmt.Printf("found mock file %s\n", entry.Name())
		r, err := readMockFile(path, entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		mocks = append(mocks, r)
	}
	return mocks, nil
}

func readMockFile(path, name string) (dulcamara.Rule, error) {
	response, err := os.ReadFile(path + "/" + name)
	if err != nil {
		return dulcamara.Rule{}, fmt.Errorf("error reading file %s:%w", path, err)
	}
	return dulcamara.Rule{Name: name, Response: string(response)}, nil
}

func main() {

	var mockFilesPath = "."
	if len(os.Args) > 1 {
		if !isDirectory(os.Args[1]) {
			fmt.Printf("Provided argument %s is not a directory. Switching to .\n", os.Args[1])
		} else {
			mockFilesPath = os.Args[1]
		}
	}
	fmt.Println("mock files path:" + mockFilesPath)
	mockFilesPath = strings.TrimSuffix(mockFilesPath, "/")
	rules, err := readMockFiles(mockFilesPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, rule := range rules {
		endpoint, err := dulcamara.ParseRule(rule)
		if err != nil {
			fmt.Println(err)
			continue
		}
		dulcamara.Deploy(endpoint)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fileEventHandler(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(mockFilesPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func fileEventHandler(event fsnotify.Event) {
	if !strings.HasSuffix(event.Name, fileSuffix) {
		return
	}
	switch {
	case event.Has(fsnotify.Write):
		fmt.Println("added/changed ", event.Name)

	case event.Has(fsnotify.Remove):
		ruleName := event.Name[strings.LastIndex(event.Name, "/")+1:]
		dulcamara.Undeploy(ruleName)
	}
}
