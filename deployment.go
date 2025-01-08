package scalingo_deployer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	scalingo "github.com/Scalingo/go-scalingo/v6"
)

var ctx context.Context

func newScalingoClient(config Config) *scalingo.Client {
	clientConfig := scalingo.ClientConfig{
		APIEndpoint: config.ScalingoAPIEndpoint,
		APIToken:    config.ScalingoAPIToken,
	}
	client, err := scalingo.New(ctx, clientConfig)
	if err != nil {
		panic(err)
	}
	return client
}

func buildDeploymentOutput(config Config, client *scalingo.Client, deployment *scalingo.Deployment) string {
	res, err := client.DeploymentLogs(ctx, deployment.Links.Output)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(
		"Deployment of %s@%s to %s finished with %s\nLog Output:\n%s\n",
		config.GitRef,
		config.GithubOwnerRepo,
		config.ScalingoApp,
		deployment.Status,
		string(body),
	)
}

func waitToFinish(config Config, client *scalingo.Client, deploymentID string) bool {
	var lastStatus scalingo.DeploymentStatus
	for {
		deployment, err := client.Deployment(ctx, config.ScalingoApp, deploymentID)
		if err != nil {
			panic(err)
		}
		if deployment.HasFailed() {
			output := buildDeploymentOutput(config, client, deployment)
			log.Print(output)
			return false
		} else if deployment.IsFinished() {
			output := buildDeploymentOutput(config, client, deployment)
			log.Print(output)
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
	ctx = context.Background()
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
	deployment, err := client.DeploymentsCreate(ctx, config.ScalingoApp, &params)
	if err != nil {
		panic(err)
	}
	if !waitToFinish(config, client, deployment.ID) {
		os.Exit(1)
	}
}
