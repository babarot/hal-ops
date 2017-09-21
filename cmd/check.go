package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/b4b4r07/hal-ops/command"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check",
	Long:  "check",
	RunE:  check,
}

func convertString(b []byte) string {
	return strings.TrimSuffix(string(b), "\n")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func check(cmd *cobra.Command, args []string) error {
	// c := command.New("bash hoge.sh")
	// if err := c.Run(); err != nil {
	// 	panic(err)
	// }
	// res := c.Result()
	// fmt.Println(res.StdoutString())
	// fmt.Printf("%#v %#v %#v\n", res.StdoutString(), res.StderrString(), res)
	// return nil

	var (
		branch string
		// stdout, stderr bytes.Buffer
	)
	if len(args) > 0 {
		branch = args[0]
		c := command.New("git checkout " + branch)
		if err := c.Run(); err != nil {
			return err
		}
		res := c.Result()
		if res.Failed {
			return fmt.Errorf(res.StderrString())
		}
	}
	fmt.Println(branch)
	return nil

	/*
		c := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		c.Stdout = &stdout
		c.Stderr = &stderr
		if err := c.Run(); err != nil {
			return err
		}
		if stderr.Len() > 0 {
			return fmt.Errorf("hoge")
		}

		branch = convertString(stdout.Bytes())
		if branch == "master" {
			return fmt.Errorf("Error: current branch is master")
		}

		// if err := cli.Run("cp", ".hal/config", filepath.Join(os.Getenv("HOME"), ".hal", "config")); err != nil {
		// 	return err
		// }

		dst := filepath.Join(os.Getenv("HOME"), ".hal", "config")
		if err := copyFile(".hal/config", dst); err != nil {
			return err
		}

		return cli.Run("hal", "config")
	*/
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
