package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

func execute(ctx context.Context, cmd string, args ...string) error {
	command := exec.CommandContext(ctx, cmd, args...)

	cmdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	cmderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	go printOutput(ctx, cmdout, os.Stdout)
	go printOutput(ctx, cmderr, os.Stderr)

	if err = command.Start(); err != nil {
		return err
	}

	err = command.Wait()
	return err
}

func printOutput(ctx context.Context, r io.Reader, output io.Writer) {
	reader := bufio.NewReader(r)

	write := func(b []byte) {
		if verbose {
			output.Write([]byte("    "))
		}

		output.Write(b)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		b, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		write(b)
	}
}

func executeString(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)

	out, err := c.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
