package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

//HandleError takes an error and does basic error handling
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

/*ScanGitFolders is a recursive function that takes a slice of strings "folders" and a string "folder",
opens the folder, and recursively scans the files in said folder, appending any .git folders to the folders slice*/
func ScanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	HandleError(err)

	files, err := f.ReadDir(-1)
	f.Close()
	HandleError(err)

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}

			folders = ScanGitFolders(folders, path)
		}
	}

	return folders
}

/*RecursiveScanFolder takes a string "folder" and calls on ScanGitFolders giving it an empty
 string slice and the folder*/
func RecursiveScanFolder(folder string) []string {
	return ScanGitFolders(make([]string, 0), folder)
}

//GetDotPath returns the file path of the file where the .git directories' file paths are being stores
func GetDotPath() string {
	usr, err := user.Current()
	HandleError(err)

	dotfile := usr.HomeDir + "/.gogitlocalstats"

	return dotfile
}

//OpenFile takes a string "filepath" and handles the opening of the file at that path, returning said file
func OpenFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filepath)
			HandleError(err)
		} else {
			panic(err)
		}
	}

	return f
}

/*JoinSlice takes two string slices "new" and "existing" and appends into "existing" the values
in "new" that don't exist in it*/
func JoinSlice(new, existing []string) []string {
	for _, i := range new {
		if !SliceContains(existing, i) {
			existing = append(existing, i)
		}
	}

	return existing
}

//SliceContains takes a string slice and a string value and checks whether that value exists in the slice
func SliceContains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

//DumpStringsToFile takes a string slice and a filepath and dumps the content of the slice into the file at the path
func DumpStringsSliceToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filepath, []byte(content), 0755)
}

/*ParseFileLinesToSlice takes a string "filepath" and parses the content
of the file at that path into a slice of strings*/
func ParseFileLinesToSlice(filepath string) []string {
	f := OpenFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return lines
}

/*AddNewSliceElementsToFile takes a string "filepath" and a slice of strings "newRepos" and calls on the functions
ParseFileLinesToSlice, JoinSlice, and DumpStringsSliceToFile*/
func AddNewSliceElementsToFile(filepath string, newRepos []string) {
	existingRepos := ParseFileLinesToSlice(filepath)
	repos := JoinSlice(newRepos, existingRepos)
	DumpStringsSliceToFile(repos, filepath)
}

//Scan takes a string "folder" and is the initial function that is used to start the scanning of the given folder
func Scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := RecursiveScanFolder(folder)
	filePath := GetDotPath()
	AddNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccefully added\n\n")
}
