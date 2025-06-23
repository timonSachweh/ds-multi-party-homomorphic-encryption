package scheduling

import (
	"log"

	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/services"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/utils"

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
	utils.PrintMemoryStats("Scheduler - updateModel")
	utils.PrintTimeForFunction("Scheduler - updateModel", mlService.RetrainAndSendUpdatedModelWeights)
}
