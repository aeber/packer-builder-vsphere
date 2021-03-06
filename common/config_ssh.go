package common

import (
	"fmt"
	packerssh "github.com/hashicorp/packer/communicator/ssh"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/multistep"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func CommHost(host string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		if host != "" {
			return host, nil
		} else {
			return state.Get("ip").(string), nil
		}
	}
}

func SshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	comm := state.Get("comm").(*communicator.Config)

	var auth []ssh.AuthMethod

	if comm.SSHPrivateKeyFile != "" {
		privateKey, err := ioutil.ReadFile(comm.SSHPrivateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Error loading configured private key file: %s", err)
		}

		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("Error setting up SSH config: %s", err)
		}

		auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		auth = []ssh.AuthMethod{
			ssh.Password(comm.SSHPassword),
			ssh.KeyboardInteractive(
				packerssh.PasswordKeyboardInteractive(comm.SSHPassword)),
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:            comm.SSHUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientConfig.Auth = auth

	return clientConfig, nil
}
