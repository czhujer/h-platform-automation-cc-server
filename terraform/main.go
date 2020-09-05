package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"log"
)

const terraformPath = "/opt/terraform_0.13.2"
const terraformWorkingDir = "/root/h-platform-automation-core/tf-owncloud"

func Run() error {
	//tmpDir, err := ioutil.TempDir("", "tfinstall")
	//if err != nil {
	//	panic(err)
	//}
	//defer os.RemoveAll(tmpDir)

	execPath, err := tfinstall.Find(tfinstall.LatestVersion(terraformPath, false))
	if err != nil {
		log.Printf("terraform: install failed: %s", err)
		return err
	}

	tf, err := tfexec.NewTerraform(terraformWorkingDir, execPath)
	if err != nil {
		log.Printf("terraform: setup error: %s", err)
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		log.Printf("terraform: init error: %s", err)
		return err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Printf("terraform: show error: %s", err)
		return err
	}

	fmt.Println(state.FormatVersion)

	return nil
}
