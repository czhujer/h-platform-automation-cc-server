package terraform

type Terraform struct {
}

//func Show() error {
//
//	const terraformPath = "/opt/terraform_0.13.2/terraform"
//	const terraformWorkingDir = "/root/h-platform-automation-core/tf-owncloud"
//
//	tf, err := tfexec.NewTerraform(terraformWorkingDir, terraformPath)
//	if err != nil {
//		log.Printf("terraform: setup error: %s", err)
//		return err
//	}
//
//	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
//	if err != nil {
//		log.Printf("terraform: init error: %s", err)
//		return err
//	}
//
//	state, err := tf.Show(context.Background())
//	if err != nil {
//		log.Printf("terraform: show error: %s", err)
//		return err
//	}
//
//	log.Printf("terraform: show: \"state.FormatVersion\": %s", state.FormatVersion)
//	log.Printf("terraform: show: \"state.Values\": %s", state.Values)
//
//	return nil
//}
