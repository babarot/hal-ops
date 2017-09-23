package cmd

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"github.com/b4b4r07/hal-ops/command/git"
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
		return fmt.Errorf("config.intergration.github_repo: should be slash char like owner/repo")
	}
	s := strings.Split(config.Conf.Integration.GitHubRepo, "/")
	owner := s[0]
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
		pull, resp, err := cli.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return err
		}
		pulls = append(pulls, pull...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	branch, err := git.GetCurrentBranchName()
	if err != nil {
		return err
	}

	num := 0
	for _, pull := range pulls {
		if pull.Head == nil {
			continue
		}
		if pull.Head.Ref == nil {
			continue
		}
		if *pull.Head.Ref == branch {
			num = pull.GetNumber()
			if num != 0 {
				break
			}
		}
	}

	if num == 0 {
		return fmt.Errorf("Not found that P-R")
	}

	msg := "Merge automatically by hal-ops"
	_, _, err = cli.PullRequests.Merge(ctx, owner, repo, num, msg, &github.PullRequestOptions{})
	if err != nil {
		return err
	}

	// Send event log to Datadog

	// Notify Slack

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
