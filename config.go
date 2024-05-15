package scalingo_deployer

import "strings"

type Config struct {
	GithubOwnerRepo     string `split_words:"true" required:"true" desc:"github repo in owner/repo format"`
	GithubAPIToken      string `split_words:"true" required:"true" desc:"github API token with read access to github repo content"`
	GitRef              string `split_words:"true" required:"true" desc:"the git reference to deploy"`
	ScalingoAPIEndpoint string `split_words:"true" required:"true" default:"https://api.osc-fr1.scalingo.com" desc:"the scalingo API endpoint to use"`
	ScalingoApp         string `split_words:"true" required:"true" desc:"the name of the scalingo app to deploy"`
	ScalingoAPIToken    string `split_words:"true" required:"true" desc:"the scalingo API token to use"`
}

func (config *Config) GithubOwner() string {
	s := strings.SplitN(config.GithubOwnerRepo, "/", 2)
	return s[0]
}

func (config *Config) GithubRepo() string {
	s := strings.SplitN(config.GithubOwnerRepo, "/", 2)
	return s[1]
}
