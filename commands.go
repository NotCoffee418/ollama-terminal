package main

import (
	"fmt"
	"os/exec"
)

func sourceBashrc() {
	cmd := exec.Command("bash", "-c", "source ~/.bashrc")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
