package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var configPath = []string{
	"$XDG_CONFIG_HOME/gopathrc",
	"$HOME/.config/gopathrc",
	"$HOME/.gopathrc",
}

type config struct {
	Command map[string]string `json:"command"`
}

func setPath() error {
	p, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	for p != "/" {
		fis, err := ioutil.ReadDir(p)
		if err != nil {
			return err
		}
		for _, fi := range fis {
			if ((fi.Mode()&os.ModeDir != 0) || (fi.Mode()&os.ModeSymlink != 0)) &&
				fi.Name() == "src" {
				return os.Setenv("GOPATH", p)
			}
		}
		p = filepath.Dir(p)
	}
	return errors.New("unable to guess GOPATH")
}

func parseConfig() *config {
	var data []byte
	for _, fname := range configPath {
		f, err := os.Open(os.ExpandEnv(fname))
		if err != nil {
			continue
		}

		data, err = ioutil.ReadAll(f)
		f.Close()
		if err == nil {
			break
		}
	}

	c := &config{}
	json.Unmarshal(data, c)

	// make a little user friendly
	if c.Command == nil {
		c.Command = make(map[string]string)
		c.Command["go"] = "/usr/bin/go"
	}
	return c
}

func main() {
	// find the original binary.
	args := os.Args
	if filepath.Base(args[0]) == "gopath" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "This program is not meant to be called directly,"+
			" please see the README for usage.\n")
		os.Exit(0)
	}

	cfg := parseConfig()
	cmd := filepath.Base(args[0])
	bin, ok := cfg.Command[cmd]
	if !ok {
		bin, _ = exec.LookPath(cmd + ".bin")
	}

	// set GOPATH if we cannot find it in environment.
	if os.Getenv("GOPATH") == "" {
		if err := setPath(); err != nil {
			fmt.Fprintf(os.Stderr, "gopath: %s\n", err)
		} else {
			// make some noise to let caller know the underlying work, for debug purpose.
			if len(args) == 1 {
				fmt.Fprintf(os.Stderr, "goapth: GOPATH set to %s\n", os.Getenv("GOPATH"))
			}
		}
	}

	c := exec.Cmd{
		Path:   bin,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := c.Run()
	if err == nil {
		os.Exit(0)
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			os.Exit(ws.ExitStatus())
		}
	}
	fmt.Fprintf(os.Stderr, "gopath: error when calling original binary(%s): %s\n",
		bin, err)
	os.Exit(1)
}
