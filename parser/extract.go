package parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
)

// ExtractAll search for compressed files and extract them in place
func ExtractAll(rootDir string) error {
	log.Printf("searching for compressed files on: %s", rootDir)
	var errors []string
	var ar archiver.Archiver
	filepath.Walk(rootDir, func(path string, file os.FileInfo, err error) (e error) {
		ar = archiver.MatchingFormat(file.Name())
		if ar == nil {
			return
		}
		if !ar.Match(path) {
			return
		}
		log.Printf("compressed file founded: %s", path)
		openErr := ar.Open(path, filepath.Dir(path))
		if openErr != nil {
			log.Printf("error extracting %s: %s", file.Name(), openErr)
			errors = append(errors, openErr.Error())
			return
		}
		log.Printf("compressed file extracted to: %s", filepath.Dir(path))
		err = os.Remove(path)
		if err != nil {
			log.Printf("error removing rar: %s", err)
		}
		log.Printf("compressed file removed: %s", file.Name())
		return nil
	})
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}
