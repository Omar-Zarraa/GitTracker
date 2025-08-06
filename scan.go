package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

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

func RecursiveScanFolder(folder string) []string {
	return ScanGitFolders(make([]string, 0), folder)
}

func GetDotPath() string {
	usr, err := user.Current()
	HandleError(err)

	dotfile := usr.HomeDir + "/.gogitlocalstats"

	return dotfile
}

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

func JoinSlice(new, existing []string) []string {
	for _, i := range new {
		if !SliceContains(existing, i) {
			existing = append(existing, i)
		}
	}

	return existing
}

func SliceContains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

func DumpStringsSliceToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filepath, []byte(content), 0755)
}

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

func AddNewSliceElementsToFile(filepath string, newRepos []string) {
	existingRepos := ParseFileLinesToSlice(filepath)
	repos := JoinSlice(newRepos, existingRepos)
	DumpStringsSliceToFile(repos, filepath)
}

func Scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := RecursiveScanFolder(folder)
	filePath := GetDotPath()
	AddNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccefully added\n\n")
}
