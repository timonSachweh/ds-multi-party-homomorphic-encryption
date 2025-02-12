package scheduling

import (
	"log"

	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/ml"
)

func InitializeCrons(mlService ml.MLService) *cron.Cron {
	c := cron.New()

	log.Println("Initializing crons")
	c.AddFunc("@every 00h20m0s", func() { updateModel(mlService) })
	c.Start()
	log.Println("Crons started")

	return c
}

func updateModel(mlService ml.MLService) {
	log.Println("Scheduler: Updating model")
	mlService.RetrainAndSendUpdatedModelWeights()
}
