package git

import (
	"fmt"

	"github.com/b4b4r07/hal-ops/command"
)

const (
	MASTER = "master"
)

func GetCurrentBranchName() (branch string, err error) {
	c := command.New("git rev-parse --abbrev-ref HEAD")
	if err := c.Run(); err != nil {
		return branch, err
	}
	branch = c.Result().StdoutString()
	return
}

func Checkout(target string) error {
	c := command.New("git checkout " + target)
	if err := c.Run(); err != nil {
		return err
	}
	res := c.Result()
	if res.Failed {
		return fmt.Errorf(res.StderrString())
	}
	return nil
}
