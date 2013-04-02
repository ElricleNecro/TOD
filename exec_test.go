package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestLaunchTask(t *testing.T) {
	command := exec.Command("tr", "a-z", "A-Z")
	command.Stdin = strings.NewReader("Un petit test")
	var out bytes.Buffer
	command.Stdout = &out

	res := LaunchTask(command, nil)
	if <-res {
		t.Log(out.String())
	}

	if out.String() != "UN PETIT TEST" {
		t.Error("Loupé, la commande a donnée : ", out.String(), " au lieu de ", "UN PETIT TEST")
	}
}
