package remote

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const defaultPrometheusServer = "192.168.121.76"
const defaultPrometheusPort = 22
const defaultPrometheusUser = "hpa-remote-executor"
const ccServerSshKey = "/root/.ssh/id_rsa"

//TODO
// remove hardcoded vmNameFull
const vmNameFull = "oc-306.hcloud.cz"

var monitoringTypes []string

type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

func loadMonitoringTypes() {
	monitoringTypes = append(monitoringTypes, "prom-file-sd-node-ocb2c")
	monitoringTypes = append(monitoringTypes, "prom-file-sd-mysql-ocb2c")
	monitoringTypes = append(monitoringTypes, "prom-file-sd-apache-ocb2c")
	monitoringTypes = append(monitoringTypes, "prom-file-sd-redis-ocb2c")
	monitoringTypes = append(monitoringTypes, "prom-file-sd-php-fpm-ocb2c")
}

func (client *SSHClient) runCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	session, err = client.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = client.prepareCommand(session, cmd)
	if err != nil {
		//log.Printf("prometheusRemote: prepareCommand failed: %s\n", err)
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		err := session.Setenv(variable[0], variable[1])
		if err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	//log.Printf("prometheusRemote: [%s:%d] creating connection\n", client.Host, client.Port)
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)

	if err != nil {
		log.Printf("prometheusRemote: [%s:%d] ssh connection failed!\n", client.Host, client.Port)
		return nil, fmt.Errorf("Failed to dial: %s", err)
	} else {
		//log.Printf("prometheusRemote: [%s:%d] ssh connection created\n", client.Host, client.Port)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func publicKeyFile(prometheusServer string, prometheusPort int, file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("prometheusRemote: [%s:%d] publicKey read failed: %s\n", prometheusServer, prometheusPort, err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Printf("prometheusRemote: [%s:%d] publicKey parse failed: %s\n", prometheusServer, prometheusPort, err)
		return nil
	}
	return ssh.PublicKeys(key)
}

func AddTarget() error {
	var err error
	var returnErr error
	var cmd *SSHCommand
	var cmdStdout bytes.Buffer
	var cmdStderr bytes.Buffer
	var cmdStdoutString string
	var cmdStderrString string

	returnErr = nil

	pubKeyRs := publicKeyFile(defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey)
	if pubKeyRs == nil {
		return fmt.Errorf("prometheusRemote: [%s:%d] unable to load publicKeyFile: %s", defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey)
	}

	sshConfig := &ssh.ClientConfig{
		User: defaultPrometheusUser,
		Auth: []ssh.AuthMethod{
			publicKeyFile(defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey),
		},
		//TODO
		// add check ssh fingerprint
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client := &SSHClient{
		Config: sshConfig,
		Host:   defaultPrometheusServer,
		Port:   defaultPrometheusPort,
	}

	loadMonitoringTypes()
	for _, monType := range monitoringTypes {
		//fmt.Println(monType)
		cmdStdout.Reset()
		cmdStderr.Reset()
		cmdStdoutString = ""
		cmdStderrString = ""

		cmd = &SSHCommand{
			Path:   fmt.Sprintf("%s %s", monType, vmNameFull),
			Stdout: &cmdStdout,
			Stderr: &cmdStderr,
		}

		log.Printf("prometheusRemote: [%s:%d] running command: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmd.Path)
		err = client.runCommand(cmd)

		cmdStdoutString = cmdStdout.String()
		cmdStderrString = cmdStderr.String()

		// serialize multiline output into one line output
		// https://github.com/bored-engineer/ssh/commit/c1e5782a7327a7b87b17dd6035df4d463dd32689
		cmdStdoutString = strings.ReplaceAll(cmdStdoutString, "\r\n", `\n`)
		cmdStderrString = strings.ReplaceAll(cmdStderrString, "\r\n", `\n`)

		log.Printf("prometheusRemote: [%s:%d] Stdout: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmdStdoutString)
		log.Printf("prometheusRemote: [%s:%d] Stderr: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmdStderrString)

		//TODO
		// return cmdStdoutString, cmdStderrString

		if err != nil {
			log.Printf("prometheusRemote: [%s:%d] command run error: %s\n", defaultPrometheusServer, defaultPrometheusPort, err)
			returnErr = err
		}
	}
	return returnErr
}

func RemoveTarget() error {
	//TODO
	// add logic

	return nil
}
