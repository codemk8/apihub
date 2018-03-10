package helm

// Run system
import (
	"fmt"
	"log"
	"os/exec"
)

// Disclaimer: This is just a quick/dirty way to communicate with tiller using helm cmd tool
// A more "programmatic" way is to use the tiller's grpc interface

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
func Install(chartName string, releaseName string) error {
	cmd := exec.Command("helm", "install", chartName, "--name="+releaseName)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error install release %s:%v:\n", "--name="+chartName, err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return err
}
