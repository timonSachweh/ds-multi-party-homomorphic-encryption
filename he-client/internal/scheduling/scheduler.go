package scheduling

import (
	"github.com/robfig/cron/v3"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/ml"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/privacy"
	"log"
)

func InitializeCrons(mlService ml.MLService, heService privacy.HEService) *cron.Cron {
	c := cron.New()

	log.Println("Initializing crons")
	c.AddFunc("@every 00h00m10s", func() { updateModel(mlService, heService) })
	c.Start()
	log.Println("Crons started")

	return c
}

func updateModel(mlService ml.MLService, heService privacy.HEService) {
	log.Println("Scheduler: Updating model")
	modelValues := mlService.GetModel().AsFloatVector()
	encryptedModel, err := heService.Encrypt(modelValues)
	if err != nil {
		log.Fatal("Scheduler: Error while encrypting model matrix")
		return
	}
	decryptedModel, err := heService.Decrypt(encryptedModel)
	if err != nil {
		log.Fatal("Scheduler: Error while decrypting model matrix")
		return
	}
	log.Println(modelValues)
	log.Println(decryptedModel)
}
