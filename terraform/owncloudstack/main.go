package owncloudstack

import (
	"fmt"
	"log"
	"os/exec"
)

func Create() error {
	//TODO
	// add tracing support
	var terraformScriptDir = "/root/h-platform-automation-core/tf-owncloud/scripts"

	cmdGenerate := exec.Command(fmt.Sprintf("%s/tf-owncloud-generate-tfvars.sh", terraformScriptDir))
	outputGenerate, err := cmdGenerate.CombinedOutput()

	//TODO
	// add terraform output to http response
	if err != nil {
		log.Printf("terraform: generate tfvars failed: %s", err)
	}
	log.Printf("terraform: generate tfvars output:")
	log.Printf("%s", outputGenerate)

	if err != nil {
		return err
	}

	cmdRun := exec.Command(fmt.Sprintf("%s/tf-owncloud-run-terraform.sh", terraformScriptDir))
	outputRun, err := cmdRun.CombinedOutput()

	//TODO
	// add terraform output to http response
	if err != nil {
		log.Printf("terraform: run failed: %s", err)
	}
	log.Printf("terraform: run output:")
	log.Printf("%s", outputRun)
	if err != nil {
		return err
	}

	return nil
}
