package ssh

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type ExecResults int

const (
	SUCCESSFUL ExecResults = iota
	FAILED
)

type SSH struct {
	hostKeyCallBack ssh.HostKeyCallback
	client          *ssh.Client
	session         ssh.Session

	Host     string
	Port     int
	Username string
	Password string
}

func New(host string, port int, username string, password string) *SSH {
	ssh := &SSH{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	return ssh
}

func (_ssh *SSH) addHostKeyCallBack() error {
	hkcb, err := knownhosts.New("C:\\Users\\Pouya\\.ssh\\known_hosts")
	if err != nil {
		return err
	}
	_ssh.hostKeyCallBack = hkcb
	return nil
}

func (_ssh *SSH) createClient() error {
	ccfg := ssh.ClientConfig{
		User: _ssh.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(_ssh.Password),
		},
		HostKeyCallback: _ssh.hostKeyCallBack,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", _ssh.Host, _ssh.Port), &ccfg)
	if err != nil {
		return err
	}
	_ssh.client = client
	return nil
}

func (_ssh *SSH) createSession() error {
	session, err := _ssh.client.NewSession()
	if err != nil {
		return err
	}
	_ssh.session = *session
	return nil
}

func (ssh *SSH) Connect() error {
	sequence := []func() error{
		ssh.addHostKeyCallBack,
		ssh.createClient,
	}
	for _, step := range sequence {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func reader(stdOut io.Reader, cancellation *bool, output chan<- string, error chan<- error) {
	buffer := make([]byte, 1024)
	for !*cancellation {
		len, err := stdOut.Read(buffer)
		if err != nil {
			if err != io.EOF {
				error <- err
			}
			break
		}
		output <- string(buffer[:len])
	}
}

func execErrorHandler(_errors <-chan string, _error <-chan error, cancellation *bool, err *error) {
	select {
	case _err := <-_error:
		{
			*cancellation = true
			*err = _err
		}
	case _err := <-_errors:
		{
			*cancellation = true
			*err = errors.New(_err)
		}
	}
}

func (ssh *SSH) Exec(command string, output chan<- string, errors chan string) error {
	if err := ssh.createSession(); err != nil {
		return err
	}
	defer ssh.session.Close()
	cancellation := false
	cmdErr := make(chan error)
	var serverErr error
	stdOut, err := ssh.session.StdoutPipe()
	if err != nil {
		return err
	}
	stdErr, err := ssh.session.StderrPipe()
	if err != nil {
		return err
	}
	go reader(stdOut, &cancellation, output, cmdErr)
	go reader(stdErr, &cancellation, errors, cmdErr)
	go execErrorHandler(errors, cmdErr, &cancellation, &serverErr)
	_, err = ssh.session.CombinedOutput(command)
	cancellation = true
	if err != nil {
		return err
	}
	if serverErr != nil {
		return serverErr
	}
	return nil
}
func (ssh *SSH) ExecMany(commands []string, output chan<- string, errors chan string) error {
	for _, command := range commands {
		if err := ssh.Exec(command, output, errors); err != nil {
			return err
		}
	}
	return nil
}
