package entities

import (
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type PrivacyParams struct {
	CKKSParameters ckks.Parameters `json:"ckks_parameters"`
	//Crs            sampling.PRNG   `json:"crs"`
}

type CkgShareExchange struct {
	Share multiparty.PublicKeyGenShare `json:"share"`
}

type PublicKeyExchange struct {
	PublicKey multiparty.PublicKeyGenShare `json:"public_key"`
}

type RelinearizationKeyShare struct {
	ShareOne multiparty.RelinearizationKeyGenShare `json:"share_one"`
	ShareTwo multiparty.RelinearizationKeyGenShare `json:"share_two"`
}
