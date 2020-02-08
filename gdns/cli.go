package gdns

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
)

func Edit() {
	f, err := ioutil.TempFile("", "gdnsctl")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	// write utf8 bom
	bom := []byte{0xef, 0xbb, 0xbf}
	_, err = f.Write(bom)
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = fmt.Fprint(f, GetHosts())
	if err != nil {
		logrus.Fatal(err)
	}

	// get os editor
	editor := "vim"
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
	raw, err := ioutil.ReadFile(f.Name())
	if err != nil {
		logrus.Fatal(err)
	}
	PutHosts(string(bytes.TrimPrefix(raw, bom)))
}

func Dump(f string) {
	if f != "" {
		err := ioutil.WriteFile(f, []byte(GetHosts()), 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		fmt.Print(GetHosts())
	}
}

func Upload(f string) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		logrus.Fatal(err)
	}
	PutHosts(string(bs))
}
