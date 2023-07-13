package main

import (
	utilities "FileManager/Classes"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
)


func getDrives() (r []string){
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
			f.Close()
		}
	}
	return 
}

func main() {

	drives := getDrives()

	fmt.Println("Hello World")
	//Create app and new window
	a := app.New()
	w := a.NewWindow("test window")
	w.Resize(fyne.NewSize(400, 400))

	//test a widget

	main_container := container.NewVBox()
	
	main_container.Add(widget.NewLabel("Hello World"))
	for _, drive := range drives {
		drive_cont := container.NewHBox()
		drive_cont.Add(widget.NewLabel("Disk " + drive))
		drive_cont.Add(widget.NewButton("Disk click",func() {
			select_drive(drive, w, "")
		} ))
		main_container.Add(drive_cont)
	}
	main_container.Add(widget.NewButton("Search", func() {
		results := deepSearch("../", "")
		fmt.Println(results)
	}))
	
	w.SetContent(main_container )
	w.ShowAndRun()
}


func select_drive(drive string, window fyne.Window, searchString string) {
	fmt.Println("Disk " + drive + " clicked")
	files, err := ioutil.ReadDir(drive + ":\\")
	if err != nil {
		fmt.Println(err)
	}
	var filenames []string
	for _, file := range files {
		filename := file.Name()
		if searchString != "" {
			if fuzzy.Match(searchString, filename) {
						filenames = append(filenames, filename)
			}
		}else{
			filenames = append(filenames, filename)
		}
		
	}
	new_list := widget.NewList(
		//length callback
		func() int {
			return len(filenames)
		},
		//create callback
		func() fyne.CanvasObject{
			return container.NewMax(widget.NewButton("Test", func(){}))
		},
		//update callback
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Button).SetText(filenames[id])
		})
	new_list.Resize(fyne.NewSize(400, 300))
	input	:= widget.NewEntry()
	input.SetText(searchString)
	
	input.OnChanged = func(s string) {
		filenames = searchCurrent(drive, s)
	}

	cont := container.NewBorder(input, nil,nil,nil, new_list)
	
	//full := container.NewMax(cont)
	//listCont := container.NewMax(new_list)
	window.SetContent(cont) 
	window.Canvas().Focus(input)
}

func searchCurrent(drive string, searchString string) []string {
	files, err := ioutil.ReadDir(drive + ":\\")
	if err != nil {
		fmt.Println(err)
	}
	var filenames []string
	for _, file := range files {
		filename := file.Name()
		if searchString != "" {
			if fuzzy.Match(searchString, filename) {
						filenames = append(filenames, filename)
			}
		}else{
			filenames = append(filenames, filename)
		}
		
	}
	return filenames
}

func deepSearch(path string, searchstring string) []utilities.SearchResult {
	var results []utilities.SearchResult
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if file.IsDir() {
			results = append(results, deepSearch(path + file.Name() + "/"  , searchstring)...)

		}else	{
			if fuzzy.Match(searchstring, file.Name()) {
				results = append(results, utilities.SearchResult{Filename: file.Name(), Filepath:path})
			}
		}
	}
	return results

}