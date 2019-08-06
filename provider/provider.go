package provider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var DefaultPath = getDefaultPath()

// Resolve resolves a filepath to a provider
// Checks the CLI arg, environment and then default path
func Resolve(providerArg string) (string, error) {
	provider := providerArg

	if provider == "" {
		provider = os.Getenv("SUMMON_PROVIDER")
	}

	if provider == "" {
		providers, _ := ioutil.ReadDir(DefaultPath)
		if len(providers) == 1 {
			provider = providers[0].Name()
		} else if len(providers) > 1 {
			return "", fmt.Errorf("More than one provider found in %s, please specify one\n", DefaultPath)
		}
	}

	provider = expandPath(provider)

	if provider == "" {
		return "", fmt.Errorf("Could not resolve a provider!")
	}

	_, err := os.Stat(provider)
	if err != nil {
		return "", err
	}

	return provider, nil
}

// Call shells out to a provider and return its output
// If call succeeds, stdout is returned with no error
// If call fails, "" is return with error containing stderr
func Call(provider, specPath string) (string, error) {
	var (
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)
	cmd := exec.Command(provider, specPath)
	cmd.Env = append(os.Environ())
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()

	if err != nil {
		errstr := err.Error()
		if stdErr.Len() > 0 {
			errstr += ": " + strings.TrimSpace(stdErr.String())
		}
		return "", fmt.Errorf(errstr)
	}

	return strings.TrimSpace(stdOut.String()), nil
}

// Given a naked filename, returns a path to executable prefixed with DefaultPath
// This is so that "./provider" will work as expected.
func expandPath(provider string) string {
	// Base returns just the last path segment.
	// If it's different, that means it's a (rel or abs) path
	if path.Base(provider) != provider {
		return provider
	}
	return path.Join(DefaultPath, provider)
}

func getDefaultPath() string {
	if runtime.GOOS == "windows" {
		// Try to get the appropriate "Program Files" directory but if one doesn't
		// exist, use a hardcoded value we think should be right.
		program_files_dir := os.Getenv("ProgramW6432")
		if program_files_dir == "" {
			program_files_dir = path.Join("C:", "Program Files")
		}

		return path.Join(program_files_dir, "Cyberark Conjur", "Summon", "Providers")
	} else {
		return "/usr/local/lib/summon"
	}
}
