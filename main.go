package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func lshCd(path string) error {
	return os.Chdir(path)
}

func lshHelp() {
	fmt.Println("Mustaq LSH")
	fmt.Println("Type program names and arguments, and hit enter")
	fmt.Println("The following are builtin")
	fmt.Println("help\tcd\texit")
}

func lshExit() {
	os.Exit(0)
}

func lshLaunch(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func lshExecute(input string) error {
	input = strings.TrimSuffix(input, "\n")

	cmd := strings.Split(input, " ")

	if cmd[0] == "cd" {
		if len(cmd) < 2 {
			return errors.New("Please enter the directory path")
		}
		return lshCd(cmd[1])
	} else if cmd[0] == "help" {
		lshHelp()
		return nil
	} else if cmd[0] == "exit" {
		lshExit()
	}

	return lshLaunch(cmd)
}

func lshLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = lshExecute(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	}
}

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	lshLoop()
}
