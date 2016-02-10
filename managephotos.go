package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Imagetosave struct {
	filename            string
	originalfilepath    string
	destinationfilepath string
}

var Images []Imagetosave
var RootPath string

func main() {

	Images = make([]Imagetosave, 10, 10)

	root, err := filepath.Abs(".")
	fmt.Println("Processing path", root)
	RootPath = root
	manageErrors(err)

	err = filepath.Walk(root, processPath)
	manageErrors(err)

	fmt.Println(Images)

	for _, image := range Images {
		err = CopyFile(image.originalfilepath+string(os.PathSeparator)+image.filename, image.destinationfilepath+string(os.PathSeparator)+image.filename)
	}
	manageErrors(err)
}

// processPath process the walk through all the file and subfolder and find all the images
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

// getFolderDestination simply open the current file ang get the metadata of creation to create the path year/month/day
func getFolderDestination(fname string) string {
	f, err := os.Open(fname)
	defer f.Close()
	manageErrors(err)

	x, err := exif.Decode(f)
	manageErrors(err)

	t, err := x.DateTime()
	manageErrors(err)

	return RootPath + string(os.PathSeparator) + t.Format("2006") + string(os.PathSeparator) + t.Format("01") + string(os.PathSeparator) + t.Format("02")
}

// manageErros, simply logs all error that are coming
func manageErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// stringInSlice checks if a string is part of a slice of string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsAnImage check that the extension of the fille is in (.jpg, .jpeg, .png and .gif)
func isAnImage(fileName string) bool {
	return stringInSlice(strings.ToLower(filepath.Ext(fileName)), []string{".jpg", ".jpeg", ".png", ".gif"})
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}

	}
	os.MkdirAll(filepath.Dir(dst), 0777)
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
