package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type Runner struct {
	script []byte
}

func (r *Runner) LoadScriptFromFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	r.script, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) LoadScriptFromCommand(command []byte) {
	r.script = command
}

func (r *Runner) Run(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/usr/bin/bash", "-c", string(r.script))
	res, err := cmd.CombinedOutput()
	fmt.Println(string(res))
	if err != nil {
		return nil, err
	}
	return res, nil
}
