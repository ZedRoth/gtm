package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
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
	// 0666 - ability for all users to read and write into it
	const newItemsPermissions = 0666

	var err error

	if len(folderPath) == 0 {
		folderPath = "."
	}

	folderPath, err = filepath.Abs(folderPath)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		os.MkdirAll(folderPath, newItemsPermissions)
	}

	workingFolderPath := filepath.Join(folderPath, workingFolderName)
	gitignoreFilePath := filepath.Join(folderPath, gitignoreFileName)

	fmt.Println(workingFolderPath)
	fmt.Println(gitignoreFilePath)

	_, err = os.Stat(gitignoreFilePath)
	if !os.IsNotExist(err) {
		err = os.Remove(gitignoreFilePath)
		if err != nil {
			panic(err)
		}
	}

	gitIgnoreFile, err := os.OpenFile(gitignoreFilePath, os.O_WRONLY|os.O_CREATE, newItemsPermissions)
	if err != nil {
		panic(err)
	}
	gitIgnoreFile.WriteString("\n" + workingFolderName + "\n")

	err = os.Mkdir(workingFolderPath, newItemsPermissions)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func main() {
	var folderPath string

	argsWithoutProgram := os.Args[1:]
	if len(argsWithoutProgram) > 0 {
		folderPath = argsWithoutProgram[0]
	}

	createWorkingFolder(folderPath)
}
