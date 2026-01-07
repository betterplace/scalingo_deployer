package scalingo_deployer

import (
	"context"
	"log"
	"net/url"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// redacted returns a URL with the query parameters removed for logging
// purposes. This prevents sensitive information like API tokens from being
// exposed in logs.
func redacted(u url.URL) string {
	u.RawQuery = ""
	return u.String()
}

// archiveDownloadURL generates the download URL for a GitHub repository
// archive. It uses the GitHub API with the provided token to get the tarball
// URL for the specified reference. The URL is logged for debugging purposes
// but with query parameters redacted.
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

	log.Printf("Determined archive download URL to be %s\n", redacted(*downloadURL))
	return downloadURL.String()
}
