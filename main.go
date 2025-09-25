package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
)

func process() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Script expects only one argument, filename with commit-msg. Got %s", os.Args)
	}
	filename := os.Args[1]

	// Get current git branch name
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchBytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Error getting git branch: %v", err)
	}

	// Extract prefix before first "/" using regex
	re := regexp.MustCompile(`^(?P<changeType>[a-z]+)/(?P<ticket>[a-zA-Z]+-\d+)`)
	matches := re.FindSubmatch(branchBytes)
	if len(matches) == 0 {
		fmt.Printf("Branch name has unexpected format (expecting <str>/<STR>-<num>), skipping...")
		return nil
	}

	change := matches[re.SubexpIndex("changeType")]
	ticket := matches[re.SubexpIndex("ticket")]

	// Read the commit message file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", filename, err)
	}

	prefix := slices.Concat(bytes.ToUpper(change), []byte(" "), bytes.ToUpper(ticket), []byte(": "))
	if !bytes.HasPrefix(content, prefix) {
		content = append(prefix, content...)
	}

	// Write the modified content back to the file
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file %s: %v", filename, err)
	}

	return nil
}

func main() {
	if err := process(); err != nil {
		fmt.Printf("Error occured: %s", err)
		os.Exit(1)
	}
}
