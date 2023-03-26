package oci

type Runner interface {
	run(cmds ...string) ([]byte, error)
}
