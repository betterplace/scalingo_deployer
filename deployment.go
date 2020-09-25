package scalingo_deployer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	scalingo "github.com/Scalingo/go-scalingo"
)

var startAt time.Time

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

func buildDeploymentOutput(client *scalingo.Client, deployment *scalingo.Deployment) string {
	res, err := client.DeploymentLogs(deployment.Links.Output)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Deployment finished with %s\nLog Output:\n%s\n", deployment.Status, string(body))
}

func waitToFinish(config Config, client *scalingo.Client, deploymentID string) bool {
	var lastStatus scalingo.DeploymentStatus
	for {
		deployment, err := client.Deployment(config.ScalingoApp, deploymentID)
		if err != nil {
			panic(err)
		}
		if deployment.HasFailed() {
			output := buildDeploymentOutput(client, deployment)
			log.Print(output)
			reportHappening(config, deployment, output, false)
			return false
		} else if deployment.IsFinished() {
			output := buildDeploymentOutput(client, deployment)
			log.Print(output)
			reportHappening(config, deployment, output, true)
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
	startAt = time.Now()
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
	if !waitToFinish(config, client, deployment.ID) {
		os.Exit(1)
	}
}
