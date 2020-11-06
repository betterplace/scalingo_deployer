package scalingo_deployer

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
)

func archiveDownloadURL(config Config) string {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: config.GithubAPIToken,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	options := github.RepositoryContentGetOptions{
		Ref: config.GitRef,
	}
	downloadURL, _, err := client.Repositories.GetArchiveLink(ctx, config.GithubOwner(), config.GithubRepo(), github.Tarball, &options, true)
	if err != nil {
		panic(err)
	}

	log.Printf("Determined archive download URL to be %s\n", downloadURL.String())
	return downloadURL.String()
}
