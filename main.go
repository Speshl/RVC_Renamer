package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	frontDirectoryPossibilities := [4]string{"./F", "./Forwards", "./Forward", "./Front"}
	rearDirectoryPossibilities := [5]string{"./B", "./Backwards", "/.Backward", "./Rear", "./R"}

	foundFrontDirectory := ""
	foundRearDirectory := ""

	for _, directory := range frontDirectoryPossibilities {
		if _, err := os.Stat(directory); !os.IsNotExist(err) {
			foundFrontDirectory = directory
			break
		}
	}

	for _, directory := range rearDirectoryPossibilities {
		if _, err := os.Stat(directory); !os.IsNotExist(err) {
			foundRearDirectory = directory
			break
		}
	}

	if foundFrontDirectory == "" || foundRearDirectory == "" {
		fmt.Printf("Found Front Directory: %s\n", foundFrontDirectory)
		fmt.Printf("Found Rear Directory: %s\n", foundRearDirectory)
		fmt.Println("ERROR: One or both directories not found.")
	}

	err := filepath.Walk(foundFrontDirectory, fileWalkFunc)
	if err != nil {
		fmt.Printf("ERROR: Not able to walk found forward directory: %s\n", err.Error())
	}

	err = filepath.Walk(foundRearDirectory, fileWalkFunc)
	if err != nil {
		fmt.Printf("ERROR: Not able to walk found forward directory: %s\n", err.Error())
	}
	fmt.Println("File Renaming Completed!!!")
}

func fileWalkFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	d := info.Sys().(*syscall.Win32FileAttributeData)
	cTime := time.Unix(0, d.CreationTime.Nanoseconds())
	//fmt.Printf("Creation Time: %+v\n", cTime)

	directory, _ := filepath.Split(path)
	newName := cTime.Format("2006-01-02T15-04-05")
	extension := filepath.Ext(path)

	newNameWithExtension := fmt.Sprintf("%s%s", newName, extension)
	newFilePath := filepath.Join(cwd, directory, newNameWithExtension)
	oldFilePath := filepath.Join(cwd, path)

	fmt.Printf("Old Name: %s\n", oldFilePath)
	fmt.Printf("New Name: %s\n", newFilePath)
	err = os.Rename(oldFilePath, newFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
