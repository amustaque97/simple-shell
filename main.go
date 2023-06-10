package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func lshLaunch(args []string, args2 []string) error {
	// todo: Better approach
	// https://stackoverflow.com/a/26541826/12902317
	r, w := io.Pipe()
	cmd := exec.Command(args[0], args[1:]...)
	cmd2 := exec.Command(args2[0], args2[1:]...)
	fmt.Println(cmd.String())
	fmt.Println(cmd2.String())

	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd2.Stdin = r
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	cmd.Start()
	cmd2.Start()

	go func() {
		defer w.Close()
		cmd.Wait()
	}()

	
	return cmd2.Wait()
}

func lshExecute(input string) error {
	input = strings.TrimSuffix(input, "\n")

	cmds := strings.Split(input, "|")

	for i := 1; i <= len(cmds)-1; i += 1 {

		cmds[i-1] = strings.Trim(cmds[i-1], " ")
		cmds[i] = strings.Trim(cmds[i], " ")

		cmd := strings.Split(cmds[i-1], " ")
		cmd2 := strings.Split(cmds[i], " ")

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

		return lshLaunch(cmd, cmd2)
	}

	return nil
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
