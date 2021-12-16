package docker

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Login(username, password, repository string) error {
	args := []string{"login", "-u", username, "-p", password}

	if private, server := isPrivateRegistry(repository); private {
		args = append(args, server)
	}

	cmd := exec.Command("docker", args...)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isPrivateRegistry(repository string) (bool, string) {
	ok, repositoryPrefix := hasPrefix(repository)

	switch {
	case !ok:
		return false, ""

	case hasPortSuffix(repositoryPrefix):
		return true, repositoryPrefix

	case looksLikeDnsName(repositoryPrefix):
		return true, repositoryPrefix

	default:
		return false, ""
	}
}

func hasPrefix(repository string) (bool, string) {
	parts := strings.Split(repository, "/")
	if len(parts) < 2 {
		return false, ""
	}

	return true, parts[0]
}

func hasPortSuffix(repositoryPrefix string) bool {
	parts := strings.Split(repositoryPrefix, ":")
	if len(parts) < 2 {
		return false
	}

	suffix := parts[len(parts)-1]
	val, err := strconv.ParseInt(suffix, 10, 32)
	if err != nil {
		return false
	}

	return suffix == strconv.FormatInt(val, 10)
}

func looksLikeDnsName(repositoryPrefix string) bool {
	parts := strings.Split(repositoryPrefix, ".")

	return len(parts) > 1
}
