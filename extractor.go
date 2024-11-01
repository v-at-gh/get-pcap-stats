package main

import (
	"fmt"
	"os"
	"os/exec"
)

func extractStats(pcapPath, statsPath string, args []string) error {
	cmdArgs := append([]string{"-r", pcapPath}, args...)

	cmd := exec.Command("tshark", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	if err := os.WriteFile(statsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output to stats file: %w", err)
	}

	return nil
}
