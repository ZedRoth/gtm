package main

import (
	"io/ioutil"
	"os"
	"os/user"
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

func createWorkingFolder() {
	const workingFolderName = ".git-issue-tracker"
	const gitIgnoreFileName = ".gitignore"

	gitIgnoreFile, err := os.OpenFile(gitIgnoreFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	gitIgnoreFile.WriteString("\n" + workingFolderName + "\n")

	// Create folder with ability to for all users to read and write into it
	err = os.Mkdir(workingFolderName, 0666)
	if err != nil {
		panic(err)
	}
}

func main() {
	createWorkingFolder()
}
