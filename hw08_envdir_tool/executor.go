package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// настраиваем переменные окружения
	for k, v := range env {
		err := os.Unsetenv(k)
		if err != nil {
			log.Fatal(err)
			return 1
		}

		if !v.NeedRemove {
			err := os.Setenv(k, v.Value)
			if err != nil {
				log.Fatal(err)
				return 1
			}
		}
	}

	// готовим команду на запуск
	var c *exec.Cmd
	switch len(cmd) {
	case 0:
		log.Println(errors.New("empty arguments"))
		return 1
	case 1:
		c = exec.Command(cmd[0]) // #nosec G204
	default:
		c = exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return 1
	}
	return 0
}
