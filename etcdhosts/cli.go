package etcdhosts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

	hosts, err := getHosts(-1)
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = fmt.Fprint(f, hosts)
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

	if hosts != string(raw) {
		p := promptx.NewDefaultPrompt(func(line []rune) error {
			if strings.ToLower(strings.TrimSpace(string(line))) != "y" && strings.ToLower(strings.TrimSpace(string(line))) != "n" {
				return errors.New("Please input 'y' or 'n'.")
			} else {
				return nil
			}
		}, "Are you sure you want to update dns records(y/n)?")

		if strings.ToLower(p.Run()) == "y" {
			putHosts(string(raw))
		} else {
			logrus.Info("dns record has not been modified.")
		}
	} else {
		logrus.Warn("dns record has not been modified.")
	}

}

func Dump(outFile string, revision int64) {
	hosts, err := getHosts(revision)
	if err != nil {
		logrus.Fatal(err)
	}
	if outFile != "" {
		err := ioutil.WriteFile(outFile, []byte(hosts), 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		fmt.Println(hosts)
	}
}

func Upload(f string) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		logrus.Fatal(err)
	}
	putHosts(string(bs))
}

func Version() {
	vl := getHostsHistory()
	cfg := &promptx.SelectConfig{
		ActiveTpl:    `»  {{ .Version | cyan }} {{ "=> Revision: " | cyan }}{{ .Revision | cyan }}`,
		InactiveTpl:  `  {{ .Version | white }} {{ "=> Revision: " | white }}{{ .Revision | white }}`,
		SelectPrompt: "Version",
		SelectedTpl:  "{{ \"» Version:\" | green }} {{ .Version }}",
		DisPlaySize:  9,
		DetailsTpl: `
--------- Etcd Hosts ----------
{{ "Hosts Version:" | faint }}	{{ .Version }}
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
