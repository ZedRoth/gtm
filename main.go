package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func outCurrentUserFolderInfo() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fileOutput, err := os.Create("test.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := fileOutput.Close()
		if err != nil {
			panic(err)
		}
	}()

	files, err := ioutil.ReadDir(user.HomeDir)
	if err != nil {
		panic(err)
	}

	_, err = fileOutput.WriteString("Current User [Home] directory is " + user.HomeDir + " and contains files\n\n")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		_, err = fileOutput.WriteString(f.Name() + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func createWorkingFolder(folderPath string) {
	const workingFolderName = ".git-task-manager"
	const gitignoreFileName = ".gitignore"

	var err error

	folderPath, err = filepath.Abs(folderPath)
	if err != nil {
		panic(err)
	}

	folderPath = strings.TrimRight(folderPath, "/") + "/"
	workingFolderPath, gitignoreFilePath := folderPath+workingFolderName, folderPath+gitignoreFileName

	fmt.Println(workingFolderPath)
	fmt.Println(gitignoreFilePath)

	gitIgnoreFile, err := os.OpenFile(gitignoreFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	gitIgnoreFile.WriteString("\n" + workingFolderPath + "\n")

	// Create folder with ability to for all users to read and write into it
	err = os.Mkdir(workingFolderPath, 0666)
	if err != nil {
		panic(err)
	}
}

func main() {
	folderPath := ""
	argsWithoutProgram := os.Args[1:]
	if len(argsWithoutProgram) > 0 {
		folderPath = argsWithoutProgram[0]
	}

	createWorkingFolder(folderPath)
}
