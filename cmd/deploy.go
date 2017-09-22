package cmd

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"github.com/b4b4r07/hal-ops/config"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "-",
	Long:  "-",
	RunE:  deploy,
}

func deploy(cmd *cobra.Command, args []string) error {
	if !strings.Contains(config.Conf.Integration.GitHubRepo, "/") {
		return fmt.Errorf("config.intergration.github_repo: should be slash char like user/repo")
	}
	s := strings.Split(config.Conf.Integration.GitHubRepo, "/")
	user := s[0]
	repo := s[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Conf.Integration.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli := github.NewClient(tc)
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 200},
		Sort:        "updated",
		Direction:   "desc",
	}

	var pulls []*github.PullRequest
	for {
		pull, resp, err := cli.PullRequests.List(ctx, user, repo, opt)
		if err != nil {
			return err
		}
		pulls = append(pulls, pull...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	branch := "test"

	for _, pull := range pulls {
		if pull.Head == nil {
			continue
		}
		if pull.Head.Ref == nil {
			continue
		}
		if *pull.Head.Ref == branch {
			n := pull.GetNumber()
			if n != 0 {
				fmt.Printf("%d", n)
				break
			}
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
