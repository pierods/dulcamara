package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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

type rule struct {
	name     string
	response string
}

func readMockFiles(path string) ([]rule, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []rule{}, err
	}
	var mocks []rule

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

func readMockFile(path, name string) (rule, error) {
	response, err := os.ReadFile(path + "/" + name)
	if err != nil {
		return rule{}, fmt.Errorf("error reading file %s:%w", path, err)
	}
	return rule{name: name, response: string(response)}, nil
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

	rules, err := readMockFiles(mockFilesPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, rule := range rules {
		parseRule(rule)
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
		fmt.Println("removed ", event.Name)
	}
}
