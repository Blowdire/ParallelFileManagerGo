package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func main() {
	root := "F:\\3d models" // Set the root directory for scanning

	var wg sync.WaitGroup
	var wgScan sync.WaitGroup      // Create a wait group to synchronize goroutines
	filePaths := make(chan string) // Create a channel to receive file paths
	dirPaths := make(chan string)

	// Launch a goroutine to traverse the filesystem and send file paths to the channel
	wgScan.Add(1)
	go func() {
		traverse(root, filePaths, dirPaths, &wgScan)
		wgScan.Wait()
		fmt.Println("Done traversing filesystem")
		close(dirPaths)
		close(filePaths)
	}()

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processFiles(filePaths)
		}()
	}
	wg.Wait()

	// Wait for all goroutines to finish
}

// traverse recursively traverses the filesystem starting from the given directory
// and sends file paths to the provided channel
func traverse(dir string, filePaths chan<- string, dirPaths chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			wg.Add(1)
			go traverse(filePath, filePaths, dirPaths, wg)
			// Recurse into subdirectories
		} else {
			filePaths <- filePath
			// Send file path to the channel
		}
	}
}

// processFiles is the function that handles the processing of file paths received from the channel
func processFiles(filePaths <-chan string) {
	for filePath := range filePaths {
		//Process the file path (e.g., perform operations on the file)
		fmt.Println("Processing file:", filePath)
	}
}
