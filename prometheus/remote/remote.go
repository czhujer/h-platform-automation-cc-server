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
	var cmdStdout bytes.Buffer
	var cmdStderr bytes.Buffer

	pubKeyRs := publicKeyFile(defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey)
	if pubKeyRs == nil {
		return fmt.Errorf("prometheusRemote: [%s:%d] unable to load publicKeyFile: %s", defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey)
	}

	sshConfig := &ssh.ClientConfig{
		User: defaultPrometheusUser,
		Auth: []ssh.AuthMethod{
			publicKeyFile(defaultPrometheusServer, defaultPrometheusPort, ccServerSshKey),
		},
		// TODO
		// add check ssh fingerprint
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client := &SSHClient{
		Config: sshConfig,
		Host:   defaultPrometheusServer,
		Port:   defaultPrometheusPort,
	}

	cmd := &SSHCommand{
		Path:   "prom-file-sd-node-ocb2c oc-101.hcloud.cz",
		Stdout: &cmdStdout,
		Stderr: &cmdStderr,
	}

	log.Printf("prometheusRemote: [%s:%d] running command: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmd.Path)
	err := client.runCommand(cmd)

	cmdStdoutString := cmdStdout.String()
	cmdStderrString := cmdStderr.String()

	// TODO
	// serialize multiline output into one line output
	// code below generates wrong order of bytes, or something like that
	//cmdStdoutString = strings.ReplaceAll(cmdStdoutString, "\n", `\n`)
	//cmdStderrString = strings.ReplaceAll(cmdStderrString, "\n", `\n`)

	log.Printf("prometheusRemote: [%s:%d] Stdout: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmdStdoutString)
	log.Printf("prometheusRemote: [%s:%d] Stderr: %s\n", defaultPrometheusServer, defaultPrometheusPort, cmdStderrString)

	if err != nil {
		log.Printf("prometheusRemote: [%s:%d] command run error: %s\n", defaultPrometheusServer, defaultPrometheusPort, err)
		return err
	} else {
		// TODO
		// return cmdStdoutString, cmdStderrString
		return nil
	}
}

func RemoveTarget() error {
	// TODO
	// add logic

	return nil
}
