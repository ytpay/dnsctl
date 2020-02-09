package etcdhosts

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/mritd/promptx"

	"github.com/sirupsen/logrus"
)

func Edit() {
	f, err := ioutil.TempFile("", "dnsctl")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = fmt.Fprint(f, getHosts())
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
	putHosts(string(raw))
}

func Dump(f string) {
	if f != "" {
		err := ioutil.WriteFile(f, []byte(getHosts()), 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		fmt.Print(getHosts())
	}
}

func Upload(f string) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		logrus.Fatal(err)
	}
	putHosts(string(bs))
}

func History() {
	vl := getHostsHistory()
	cfg := &promptx.SelectConfig{
		ActiveTpl:    `»  {{ .Version | cyan }} {{ "=> Revision: " | cyan }}{{ .Revision | cyan }}`,
		InactiveTpl:  `  {{ .Version | white }} {{ "=> Revision: " | white }}{{ .Revision | white }}`,
		SelectPrompt: "Version",
		SelectedTpl:  "{{ \"» Version:\" | green }} {{ .Version }}",
		DisPlaySize:  9,
		DetailsTpl: `
--------- Etcd Hosts ----------
{{ "Version:" | faint }}	{{ .Version }}
{{ "Etcd Revision:" | faint }}	{{ .Revision }}`,
	}

	s := &promptx.Select{
		Items:  vl,
		Config: cfg,
	}
	idx := s.Run()

	f, err := ioutil.TempFile("", "dnsctl")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = fmt.Fprint(f, vl[idx].Hosts)
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
}
