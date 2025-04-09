package entities

import (
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type PrivacyParams struct {
	CKKSParameters ckks.Parameters `json:"ckks_parameters"`
	//Crs            sampling.PRNG   `json:"crs"`
	Tpk rlwe.PublicKey `json:"tpk"`
}

type CkgShareExchange struct {
	Share multiparty.PublicKeyGenShare `json:"share"`
}

type PublicKeyExchange struct {
	PublicKey rlwe.PublicKey `json:"public_key"`
}

type RelinearizationKeyShare struct {
	ShareOne multiparty.RelinearizationKeyGenShare `json:"share_one"`
	ShareTwo multiparty.RelinearizationKeyGenShare `json:"share_two"`
}

type PublicKeySwitchGeneration struct {
	EncryptedWeights [][]byte       `json:"encrypted_weights"`
	TargetPublicKey  rlwe.PublicKey `json:"target_public_key"`
}

func (p *PublicKeySwitchGeneration) WeightsAsCiphertext() []*rlwe.Ciphertext {
	ciphertexts := make([]*rlwe.Ciphertext, len(p.EncryptedWeights))
	for i, w := range p.EncryptedWeights {
		ciphertext := rlwe.Ciphertext{}
		ciphertext.UnmarshalBinary(w)
		ciphertexts[i] = &ciphertext
	}
	return ciphertexts
}

type PublicKeySwitchShare struct {
	Share multiparty.PublicKeySwitchShare `json:"share"`
}
