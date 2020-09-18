package scalingo_deployer

import (
	"log"
	"os"
	"time"

	scalingo "github.com/Scalingo/go-scalingo"
)

func newScalingoClient(config Config) *scalingo.Client {
	clientConfig := scalingo.ClientConfig{
		APIEndpoint: config.ScalingoAPIEndpoint,
		APIToken:    config.ScalingoAPIToken,
	}
	client, err := scalingo.New(clientConfig)
	if err != nil {
		panic(err)
	}
	return client
}

func waitToFinish(client *scalingo.Client, scalingoApp, deploymentID string) bool {
	var lastStatus scalingo.DeploymentStatus
	for {
		deployment, err := client.Deployment(scalingoApp, deploymentID)
		if err != nil {
			panic(err)
		}
		if deployment.HasFailed() {
			log.Printf("Deployment failed with %s\n", deployment.Status)
			return false
		} else if deployment.IsFinished() {
			log.Printf("Deployment has finished in state %s\n", deployment.Status)
			break
		} else {
			if deployment.Status != lastStatus {
				log.Printf("Deployment is running in state %s\n", deployment.Status)
			}
			time.Sleep(time.Second)
		}
		lastStatus = deployment.Status
	}
	return true
}

func Start(config Config) {
	log.Printf(
		"Starting deployment of %s@%s to scalingo app %s\n",
		config.GitRef,
		config.GithubOwnerRepo,
		config.ScalingoApp,
	)
	sourceURL := archiveDownloadURL(config)
	client := newScalingoClient(config)
	params := scalingo.DeploymentsCreateParams{
		GitRef:    &config.GitRef,
		SourceURL: sourceURL,
	}
	deployment, err := client.DeploymentsCreate(config.ScalingoApp, &params)
	if err != nil {
		panic(err)
	}
	if !waitToFinish(client, config.ScalingoApp, deployment.ID) {
		os.Exit(1)
	}
}
