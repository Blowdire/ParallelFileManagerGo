package FileManager

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func main() {
	fmt.Println("Hello World")
	path := "F:\\3d models"
	//channel used to pass results from threads
	//pathsFound := traverse(path)
	dirs := make(chan string)
	files := make(chan string)
	files_found, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	for _, file := range files_found {
		filePath := filepath.Join(path, file.Name())

		if file.IsDir() {
			dirs <- filePath // Recurse into subdirectories
		} else {
			files <- filePath
		}
	}
	var wg sync.WaitGroup
	for dir := range dirs {
		fmt.Println(dir)
		wg.Add(1)
		go traverse(dir, dirs, files, &wg)
	}

	wg.Wait()

	close(dirs)
	close(files)
	for msg := range files {
		fmt.Println(msg)
	}

}
func traverse(dir string, dirChan chan<- string, fileChan chan<- string, wg *sync.WaitGroup) {
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
			dirChan <- filePath // Recurse into subdirectories
		} else {
			fileChan <- filePath
		}
	}

}
