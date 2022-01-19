package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
)

func removeFilesFromList(files []os.FileInfo) []os.FileInfo {
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() == false {
			files = append(files[:i], files[i+1:]...)
		}
	}

	return files
}

func formatSize(file os.FileInfo) string {
	file_size := file.Size()
	size_str := fmt.Sprintf("%d", file_size)
	if size_str == "0" {
		size_str = "empty"
	} else {
		size_str += "b"
	}
	size_str = fmt.Sprintf("(%s)", size_str)
	return size_str
}

func isExecutable(mode fs.FileMode) bool {
	modeString := mode.String()
	for _, access := range modeString {
		if access == 'x' {
			return true
		}
	}
	return false
}

func printLevel(out io.Writer, path string, prepre string, prefix string, printfile bool) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	files, err := ioutil.ReadDir(file.Name())
	if err != nil {
		fmt.Println(err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	var end string
	if len(files) == 1 {
		end = "└───"
	} else {
		end = "├───"
	}

	if printfile == false {
		files = removeFilesFromList(files)
	}

	prefix = prepre
	for idx, current_file := range files {
		if idx == len(files)-1 {
			end = "└───"
		}
		if current_file.IsDir() == false {
			size_str := formatSize(current_file)
			if isExecutable(current_file.Mode().Perm()) {
				blue := color.New(color.FgGreen).SprintFunc()
				line := prefix + end + blue(current_file.Name())
				fmt.Fprintln(out, line)
			} else {
				fmt.Fprintln(out, prefix+end+current_file.Name(), size_str)
			}

		} else {
			blue := color.New(color.FgBlue).SprintFunc()
			line := prefix + end + blue(current_file.Name())
			fmt.Fprintln(out, line)
		}

		var prepre string
		if current_file.IsDir() == true {
			if idx == len(files)-1 {
				prepre = prefix + prepre + "	"
			} else {
				prepre = prefix + prepre + "│	"
			}
			printLevel(out, filepath.Join(file.Name(), current_file.Name()), prepre, prefix, printfile)
		}
	}
}

func dirTree(out io.Writer, path string, printFile bool) error {
	printLevel(out, path, "", "", printFile)
	return nil
}

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
