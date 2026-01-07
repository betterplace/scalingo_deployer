package scalingo_deployer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	scalingo "github.com/Scalingo/go-scalingo/v7"
)

var ctx context.Context

// newScalingoClient creates a new Scalingo API client with the provided
// configuration. It initializes the client with the API endpoint and
// authentication token from the config. The function panics if client creation
// fails, as this is a critical initialization step.
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

// buildDeploymentOutput retrieves and formats the deployment logs for display.
// It fetches the deployment output from the Scalingo API and formats it with
// deployment details. The function panics if it fails to retrieve the logs or
// read the response body.
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

// waitToFinish monitors a deployment until it completes or fails
// It polls the Scalingo API every second to check the deployment status
// Returns true if deployment succeeds, false if it fails
// The function panics if it encounters an error while checking deployment status
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

// Start initiates the deployment process for a GitHub repository to Scalingo.
// It performs the following steps:
// 1. Creates a Scalingo client
// 2. Determines the archive download URL from GitHub
// 3. Creates a deployment with the specified parameters
// 4. Waits for the deployment to complete
// 5. Exits with error code 1 if deployment fails
// The function panics if any step in the deployment process fails
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
