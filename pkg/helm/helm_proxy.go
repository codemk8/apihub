package helm

// Run system
import (
	"fmt"
	"log"
	"os/exec"
)

// ListRelease list the running chart releases
func ListRelease() error {
	cmd := exec.Command("helm", "list")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return err
}

// Repo runs "helm repo ..." command
func Repo(subcmd string) error {
	cmd := exec.Command("helm", "repo", subcmd)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return err
}

// Install runs "helm install ..." command
func Install(chartName string) error {
	cmd := exec.Command("helm", "install", chartName)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return err
}
