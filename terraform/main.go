package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
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
		panic(err)
	}

	tf, err := tfexec.NewTerraform(terraformWorkingDir, execPath)
	if err != nil {
		panic(err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		panic(err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(state.FormatVersion)

	return nil
}
