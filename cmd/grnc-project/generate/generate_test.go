package generate

import "testing"

func TestExit(t *testing.T) {

	allowExit = false

	g := new(ProjectGenerator)

	g.exitError("Exit %s", "message")

}
