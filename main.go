package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/frankywahl/allowedSignatures/internal/github"
	"github.com/frankywahl/allowedSignatures/internal/ssh"
)

var (
	ghToken         string
	githubEnteprise = "http://api.github.com/graphql"
	verbose         bool
	owner, repo     string
	useContributors bool
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	if err := parseFlags(ctx); err != nil {
		return err
	}

	opts := []github.Option{}

	if verbose {
		opts = append(opts, github.SetVerbose())
	}

	ghClient, err := github.NewEnterpriseClient(githubEnteprise, ghToken, opts...)
	if err != nil {
		return err
	}

	var users []github.User
	if useContributors {
		users, err = ghClient.GetContributorKeys(ctx, owner, repo)
	} else {
		users, err = ghClient.GetCollaboratorKeys(ctx, owner, repo)
	}

	if err != nil {
		return err
	}

	if err := printOutput(os.Stdout, users); err != nil {
		return err
	}

	return nil
}

func printOutput(w io.Writer, users []github.User) error {
	for _, user := range users {
		for _, key := range ssh.FilterSigningKeys(user.Keys) {
			fmt.Fprintf(os.Stdout, "%s %s %s\n", user.Login, key, user.Login)
		}
	}

	return nil
}

func parseFlags(ctx context.Context) error {
	flag.StringVar(&ghToken, "github-token", os.Getenv("GITHUB_TOKEN"), "the github token to use to make requests\ndefaults to environment variable GITHUB_TOKEN")
	flag.BoolVar(&verbose, "verbose", false, "print debugging information")
	flag.BoolVar(&useContributors, "use-contributors", false, "use contributors to generate list. This is more complete, but will make many more requests to GitHub")
	flag.StringVar(&repo, "repository", "", "the repository to get the information for")
	flag.StringVar(&githubEnteprise, "github-enterprise", "https://api.github.com/graphql", "use github enterprise URL as the endpoint instead of standard GitHub")
	flag.StringVar(&owner, "owner", "", "the organisation or owner of the repository")
	flag.Parse()

	if owner == "" {
		return errors.New("owner cannot be blank")
	}

	if repo == "" {
		return errors.New("repository cannot be blank")
	}

	return nil
}
