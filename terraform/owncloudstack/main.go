package owncloudstack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"log"
)

// TODO
// generate oc-xyz.tfvars
//
// example:
// vmname = "oc-xyz"
// vmip = "10.1.2.4"
// vm_data_disk_size = "10"

func Create() error {

	var terraformPath = "/opt/terraform_0.13.2/terraform"
	var terraformWorkingDir = "/root/h-platform-automation-core/tf-owncloud"

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
	// add PlanOptions for oc-xyz.tfvars

	rs, err := tf.Plan(context.Background())
	if err != nil {
		log.Printf("terraform: plan error: %s", err)
		return err
	}

	// TODO
	// add terraform apply

	fmt.Println(rs)

	return nil
}
