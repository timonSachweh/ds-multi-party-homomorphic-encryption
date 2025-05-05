package entities

import (
	"github.com/Pro7ech/lattigo/he/hefloat"
	"github.com/Pro7ech/lattigo/mhe"
	"github.com/Pro7ech/lattigo/rlwe"
)

type PrivacyParams struct {
	CKKSParameters hefloat.Parameters
	Tpk            rlwe.PublicKey
}

type CkgShareExchange struct {
	Share mhe.PublicKeyShare `json:"share"`
}

type PublicKeyExchange struct {
	PublicKey rlwe.PublicKey `json:"public_key"`
}

type RelinearizationKeyShare struct {
	Share mhe.RelinearizationKeyShare `json:"share"`
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
	Share mhe.KeySwitchingShare `json:"share"`
}
