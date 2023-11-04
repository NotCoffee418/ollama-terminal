package main

import (
	"fmt"
	"os/exec"
	"os/user"
)

func sudoCheck() {
	// Get the current user running the process.
	currentUser, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Unable to get the current user: %s", err))
	}

	// Check if the effective user ID (eUID) is 0 which is the superuser ID
	if currentUser.Uid == "0" {
		panic("This program should not be run as superuser (sudo)")
	}
}

func execCheck(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

func installCheck() {
	settings := LoadSettings()
	if settings.isUndefined {
		install()
	}

}
