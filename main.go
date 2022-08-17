package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed output.txt
var content string
var outstr string
var out []byte
var err error

func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func compare() {

	outstr := string(out)

	if outstr == content {

		fmt.Println("nothing changed")
	}

}

func callHR() {

	cmd := exec.Command("bash", "-c", "./handelsregister-main/run.sh")
	out, err = cmd.Output()

	if err != nil {
		fmt.Println(string(out))
		println(err.Error())
		return

	}
}

func main() {

	callHR()
	compare()
	rm := ClearDir("./cache/")
	if rm != nil {
		fmt.Println(rm)
	}

}
