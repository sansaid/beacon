//go:generate mockgen -source $GOFILE -package oci -destination ./runner_mock.go Runner
package oci

type Runner interface {
	run(cmds ...string) ([]byte, error)
}
