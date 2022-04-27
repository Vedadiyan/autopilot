package configurator

import "autopilot/internal/ssh"

func InstallGit(ssh ssh.SSH, output chan<- string, errors chan string) error {
	return ssh.Exec("dnf install -y git", output, errors)
}
