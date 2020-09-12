package owncloudstack

import (
	"fmt"
	"log"
	"os/exec"
)

func Create() error {
	var terraformScriptDir = "/root/h-platform-automation-core/tf-owncloud/scripts"

	cmd := exec.Command(fmt.Sprintf("%s/tf-owncloud-generate-tfvars.sh", terraformScriptDir))
	if err := cmd.Run(); err != nil {
		log.Printf("terraform: generate tfvars failed: %s", err)
		return err
	}

	//TODO
	// run terraform
	// wrapper script: tf-owncloud-run-terraform.sh

	return nil
}
