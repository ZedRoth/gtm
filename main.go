package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
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
	const workingFolderNamePattern = `(?m:^(\n|\r|)\.git-task-manager(\n|\r|)$)`
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

	var gitignoreContent string
	alreadyIgnored := false
	if _, err = os.Stat(gitignoreFilePath); !os.IsNotExist(err) {
		gitignoreContentBytes, err := ioutil.ReadFile(gitignoreFilePath)
		if err != nil {
			panic(err)
		}

		gitignoreContent = string(gitignoreContentBytes)
		alreadyIgnored, err = regexp.MatchString(workingFolderNamePattern, gitignoreContent)
		if err != nil {
			panic(err)
		}
	}

	if !alreadyIgnored {
		gitignoreFile, err := os.OpenFile(gitignoreFilePath, os.O_APPEND|os.O_CREATE, newItemsPermissions)
		if err != nil {
			panic(err)
		}
		defer gitignoreFile.Close()

		linesToAdd := "# Working folder of git-task-manager\r\n" + workingFolderName + "\r\n"
		if len(gitignoreContent) > 0 {
			linesToAdd = "\r\n\r\n" + linesToAdd
		}
		gitignoreFile.WriteString(linesToAdd)
	}

	err = os.Mkdir(workingFolderPath, newItemsPermissions)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

func createRepositoryIfNeed(folderPath string) {
	if len(folderPath) == 0 {
		folderPath = "."
	}

	folderPath, err := filepath.Abs(folderPath)
	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(folderPath); os.IsNotExist(err) {
		panic(err)
	}

	gitCommand := "git"
	statusArguments := []string{"status"}
	initArguments := []string{"init"}

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(folderPath)
	if err != nil {
		panic(err)
	}
	defer os.Chdir(currentWorkingDir)

	_, err = exec.Command(gitCommand, statusArguments...).Output()
	if err != nil {
		_, err = exec.Command(gitCommand, initArguments...).Output()
	}

	if err != nil {
		panic(err)
	}
}

func createTaskManagerBranchIfNeed(folderPath string) {
	const branchName = "git-task-manager"

	if len(folderPath) == 0 {
		folderPath = "."
	}

	folderPath, err := filepath.Abs(folderPath)
	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(folderPath); os.IsNotExist(err) {
		panic(err)
	}

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(folderPath)
	if err != nil {
		panic(err)
	}
	defer os.Chdir(currentWorkingDir)

	gitCommand := "git"

	revParseOutput, err := exec.Command(gitCommand, []string{"rev-parse", branchName}...).Output()
	if err != nil {
		fmt.Println("there is no such object in git repository so we are creating new one")

		localBranchesOutput, err := exec.Command(gitCommand, []string{"rev-parse", branchName}...).Output()

	} else {
		fmt.Printf("there is such object in git repository: %s\n", revParseOutput)
	}
}

func main() {
	var folderPath string

	argsWithoutProgram := os.Args[1:]
	if len(argsWithoutProgram) > 0 {
		folderPath = argsWithoutProgram[0]
	}

	createWorkingFolder(folderPath)
	createRepositoryIfNeed(folderPath)
	createTaskManagerBranchIfNeed(folderPath)
}
