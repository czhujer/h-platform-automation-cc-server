package owncloudstackdocker

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"log"
)

func Create() error {
	var err error
	var rsPlan bool
	var terraformPath = "/opt/terraform_0.13.2/terraform"
	var terraformWorkingDir = "/root/h-platform-automation-core/tf-owncloud"

	//TODO
	// add tracing support

	tf, err := tfexec.NewTerraform(terraformWorkingDir, terraformPath)
	if err != nil {
		log.Printf("terraform: setup error: %s", err)
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		log.Printf("terraform: init error: %s", err)
		return err
	}

	// TODO
	// add PlanOptions for oc-1xyz.tfvars

	rsPlan, err = tf.Plan(context.Background())
	if err != nil {
		log.Printf("terraform: plan error: %s", err)
		return err
	}

	log.Printf("terraform: apply results: %b", rsPlan)

	err = tf.Apply(context.Background())
	if err != nil {
		log.Printf("terraform: apply error: %s", err)
		return err
	}

	return nil
}
