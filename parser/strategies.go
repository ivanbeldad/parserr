package parser

import (
	"log"
	"os"
	"parserr/api"
	"path"
	"path/filepath"
)

// FixStrategy ...
type FixStrategy interface {
	Fix(m *api.Media) error
}

// MaintainPathStrategy Rename file in place if its inside a folder or
// create a folder with the name of the file and move it to that folder
type MaintainPathStrategy struct {
	API   api.RRAPI
	Mover Mover
}

// MoveToOwnFolderStrategy Move file
type MoveToOwnFolderStrategy struct {
	API   api.RRAPI
	Mover Mover
}

// StrategyFactory Return the fix strategy depending on the api
func StrategyFactory(a api.RRAPI, m Mover) FixStrategy {
	if a.GetType() == api.TypeMovie {
		return MaintainPathStrategy{
			API:   a,
			Mover: m,
		}
	}
	return MoveToOwnFolderStrategy{
		API:   a,
		Mover: m,
	}
}

// Fix Rename file in place if its inside a folder or
// create a folder with the name of the file and move it to that folder
func (s MaintainPathStrategy) Fix(m *api.Media) (err error) {
	err = s.move(m)
	if err != nil {
		return
	}
	return nil
}

func (s MaintainPathStrategy) move(m *api.Media) (err error) {
	log.Printf("fixing: %s", m.FilenameOri)
	fileLocation := m.FileLocOri
	fileIsOnRoot := m.QueueElem.Title == m.FilenameOri
	if fileIsOnRoot {
		fileLocation, err = moveFileToFolderWithSameName(m.FileLocOri, s.Mover)
		if err != nil {
			log.Printf("cannot move file to a folder: %s", err.Error())
			return err
		}
	}
	newFileLocation := path.Join(filepath.Dir(fileLocation), m.FilenameFinal)
	log.Printf("moving from %s to %s", fileLocation, newFileLocation)
	err = s.Mover.Move(fileLocation, newFileLocation)
	if err != nil {
		return err
	}
	m.FileLocFinal = newFileLocation
	return nil
}

func moveFileToFolderWithSameName(fileLocation string, m Mover) (dest string, err error) {
	log.Printf("moving file to a folder with its own name")
	tmpPath := fileLocation + ".tmp"
	err = m.Move(fileLocation, tmpPath)
	if err != nil {
		return "", err
	}
	err = m.Mkdir(fileLocation)
	if err != nil {
		m.Move(tmpPath, fileLocation)
		return "", err
	}
	dest = path.Join(fileLocation, filepath.Base(fileLocation))
	err = m.Move(tmpPath, dest)
	if err != nil {
		return
	}
	return
}

// Fix Rename file in place if its inside a folder or
// create a folder with the name of the file and move it to that folder
func (s MoveToOwnFolderStrategy) Fix(m *api.Media) (err error) {
	log.Printf("move to own folder strategy: %s", m.FilenameOri)
	err = s.moveToFolder(m)
	if err != nil {
		return
	}
	newDir := filepath.Dir(m.FileLocFinal)
	s.orderToImportFiles(newDir)
	if _, err := os.Stat(newDir); err == nil {
		log.Printf("file not imported correctly: %s", m.FileLocFinal)
		err = s.Mover.Move(m.FileLocFinal, m.FileLocOri)
		log.Printf("moving file back from: %s to: %s", m.FileLocFinal, m.FileLocOri)
		os.Remove(newDir)
		m.FileLocFinal = m.FileLocOri
	}
	return nil
}

func (s MoveToOwnFolderStrategy) moveToFolder(m *api.Media) (err error) {
	destDir := path.Join(s.API.GetDownloadFolder(), m.FilenameFinal)
	destFile := path.Join(destDir, destDir+m.FileExtension)
	s.Mover.Mkdir(destDir)
	err = s.Mover.Move(m.FileLocOri, destFile)
	if err != nil {
		log.Printf("cannot move file: %s", err.Error())
		return
	}
	m.FileLocFinal = destFile
	log.Printf("file moved, new destination: %s", m.FileLocFinal)
	return
}

func (s MoveToOwnFolderStrategy) orderToImportFiles(path string) (err error) {
	log.Printf("forcing to import files from: %s", path)
	command := s.API.DownloadScan(path)
	_, err = s.API.ExecuteCommandAndWait(command, api.DefaultRetries)
	return
}
