package scheduling

import (
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"log"

	"github.com/robfig/cron/v3"
)

func InitializeCrons(mlService services.MLService) *cron.Cron {
	c := cron.New()

	log.Println("Initializing crons")
	c.AddFunc("@every 00h10m0s", func() { updateModel(mlService) })
	c.Start()
	log.Println("Crons started")

	return c
}

func updateModel(mlService services.MLService) {
	log.Println("Scheduler: Updating model")
	mlService.RetrainAndSendUpdatedModelWeights()
}
