package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

const (
	file_print_pattern       = `%s (%sb)`
	empty_file_print_pattern = `%s (empty)`
	folder_print_pattern     = `%s`
)

const (
	empty_arrow = "\t"
	mid_arrow   = "│\t"
	in_arrow    = "├───"
	end_arrow   = "└───"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := buildTree(out, path, printFiles, "")

	if err != nil {
		panic(err.Error())
	}

	return nil
}

func buildTree(out io.Writer, path string, printFiles bool, tail string) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	realFiles := []fs.FileInfo{}
	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				realFiles = append(realFiles, file)
			}
		}
	} else {
		realFiles = files
	}

	for counter, file := range realFiles {
		newtail := tail
		newpath := path + "\\" + file.Name()
		if counter == (len(realFiles) - 1) {
			fmt.Fprintf(out, "%s%s", tail+end_arrow, file.Name())
			newtail = newtail + empty_arrow
		} else {
			fmt.Fprintf(out, "%s%s", tail+in_arrow, file.Name())
			newtail = newtail + mid_arrow
		}
		if printFiles && !file.IsDir() {
			fileSize := "empty"
			if 0 < file.Size() {
				fileSize = fmt.Sprintf("%vb", file.Size())
			}
			fmt.Fprintf(out, " (%s)\n", fileSize)
		} else {
			fmt.Fprintf(out, "\n")
		}
		buildTree(out, newpath, printFiles, newtail)
	}
	return nil
}
