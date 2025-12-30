package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

type passwordCommandError string

const (
	errPasswordCommandEmpty    passwordCommandError = "password command is empty"
	errPasswordCommandNoOutput passwordCommandError = "password command produced no password"
)

func (err passwordCommandError) Error() string {
	return string(err)
}

// PreparePassword stores the output from --password-command for storage setup.
func PreparePassword(ctx *cli.Context) error {
	command := ctx.String(passwordCommandFlagName)
	if command == "" {
		return nil
	}

	password, err := readPasswordCommand(command)
	if err != nil {
		return err
	}

	if err = os.Setenv("TOTP_PASS", password); err != nil {
		return fmt.Errorf("set password environment: %w", err)
	}

	return nil
}

func readPasswordCommand(command string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", errPasswordCommandEmpty
	}

	cmd := shellCommand(command)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("password command failed: %w", err)
	}

	password := strings.TrimRight(string(output), "\r\n")
	if password == "" {
		return "", errPasswordCommandNoOutput
	}

	return password, nil
}

func shellCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(context.Background(), "cmd", "/C", command)
	}

	return exec.CommandContext(context.Background(), "sh", "-c", command)
}
