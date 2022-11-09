package scalingo_deployer

import (
	"fmt"
	"log"
	"time"

	scalingo "github.com/Scalingo/go-scalingo/v4"
	h "github.com/flori/happening"
)

func reportHappening(config Config, deployment *scalingo.Deployment, output string, success bool) {
	if config.HappeningURL == "" {
		return
	}
	hc := h.NewConfig()
	hc.Name = fmt.Sprintf("Deployment %s", config.ScalingoApp)
	hc.ReportURL = config.HappeningURL
	event := h.Execute(*hc, nil)
	event.Started = *deployment.CreatedAt
	event.Duration = time.Now().Sub(startAt)
	event.Output = output
	event.Store = true
	if !success {
		event.ExitCode = 1
	}
	event.Success = success
	log.Printf("Sending \"%s\" to happening: %v\n", event.Name, string(h.EventToJSON(event)))
	h.SendEvent(event, hc)
}
