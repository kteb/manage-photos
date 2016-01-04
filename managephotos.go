package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"os"
	"path/filepath"
	//"time"
)

type Imagetosave struct {
	filename            string
	originalfilepath    string
	destinationfilepath string
}

var Images []Imagetosave

func main() {

	Images = make([]Imagetosave, 10, 10)

	root, err := filepath.Abs(".")
	fmt.Println("Processing path", root)
	manageErrors(err)

	err = filepath.Walk(root, processPath)
	manageErrors(err)

	fmt.Println(Images)

}

func processPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if path != "." {
		if !info.IsDir() {
			if isAnImage(info.Name()) {
				Images = append(Images, Imagetosave{info.Name(), filepath.Dir(path), getFolderDestination(path)})
			}
		}
	}

	return nil
}

func getFolderDestination(fname string) string {
	f, err := os.Open(fname)
	defer f.Close()
	manageErrors(err)

	x, err := exif.Decode(f)
	manageErrors(err)

	t, err := x.DateTime()
	manageErrors(err)

	return t.Format("2006") + string(os.PathSeparator) + t.Format("01") + string(os.PathSeparator) + t.Format("02")
}

func manageErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isAnImage(fileName string) bool {
	return stringInSlice(filepath.Ext(fileName), []string{".jpg", ".jpeg", ".png", ".gif"})
}
