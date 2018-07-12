package parser

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sonarr-parser-helper/api"
	"strings"
)

// FixFileNames Try to rename downloaded files to the original torrent name
func FixFileNames(failedFiles []*api.Media, m Mover, downloadFolder string) error {
	var errors []string
	for _, file := range failedFiles {
		err := fixFileName(file, m, downloadFolder)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, ", "))
}

func fixFileName(failedFile *api.Media, m Mover, downloadFolder string) error {
	log.Printf("fixing: %s", failedFile.QueueElement.Title)
	filename, err := failedFile.GuessFileName()
	if err != nil {
		log.Printf("impossible to guess original filename: %s", failedFile.QueueElement.Title)
		return err
	}
	finalName, err := failedFile.GuessFinalName(filename)
	if err != nil {
		log.Printf("impossible to guess a final filename: %s", failedFile.QueueElement.Title)
		return err
	}
	oldPath, err := locationOfFile(downloadFolder, filename)
	if err != nil {
		log.Printf("impossible to get location of file: %s", failedFile.QueueElement.Title)
		return err
	}
	if fileIsOnRoot(failedFile, filename) {
		oldPath, err = moveFileToFolder(oldPath, m)
		if err != nil {
			log.Printf("cannot move file to a folder")
			return err
		}
	}
	newPath := path.Join(filepath.Dir(oldPath), finalName+filepath.Ext(oldPath))
	log.Printf("renaming %s to %s", oldPath, newPath)
	err = m.Move(oldPath, newPath)
	if err != nil {
		return err
	}
	failedFile.HasBeenRenamed = true
	return nil
}

func fileIsOnRoot(failedFile *api.Media, filename string) bool {
	return failedFile.QueueElement.Title == filename
}

func moveFileToFolder(oldpath string, m Mover) (dest string, err error) {
	log.Printf("moving file to a folder with its own name")
	tmpPath := oldpath + ".tmp"
	err = m.Move(oldpath, tmpPath)
	if err != nil {
		return "", err
	}
	err = os.Mkdir(oldpath, 0775)
	if err != nil {
		m.Move(tmpPath, oldpath)
		return "", err
	}
	dest = path.Join(oldpath, filepath.Base(oldpath))
	err = m.Move(tmpPath, dest)
	if err != nil {
		return
	}
	return
}

// locationOfFile Search recursively on root for a file with filename
// and return its complete path
func locationOfFile(root, filename string) (location string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Name() == filename {
			location = path
			return fmt.Errorf("ok")
		}
		return nil
	})
	if err != nil && err.Error() == "ok" {
		err = nil
	}
	if location == "" {
		err = fmt.Errorf("%s doesn't exists inside %s", filename, root)
	}
	return
}
