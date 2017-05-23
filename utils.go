package ghpages

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Info(msg interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", msg)
}

func Prompt(msg interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", msg)
}

func CheckIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFilesList(dirPath string) []string {
	files := []string{}
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == dirPath {
			return nil
		}
		files = append(files, path)
		return nil
	})
	CheckIfErr(err)
	return files
}

func CopyFile(src string, dest string, overwrite bool) error {
	in, err := os.Open(src)
	CheckIfErr(err)
	defer in.Close()

	if stat, err := in.Stat(); stat.IsDir() {
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			err = os.Mkdir(dest, 0777)
			CheckIfErr(err)
		}
		CheckIfErr(err)
		return nil
	}

	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if !overwrite {
			return nil
		}
	}

	out, err := os.Create(dest)
	CheckIfErr(err)
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)
	CheckIfErr(err)
	return out.Sync()
}
