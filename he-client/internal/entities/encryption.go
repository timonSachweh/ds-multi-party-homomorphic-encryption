package entities

import (
	"github.com/Pro7ech/lattigo/he/hefloat"
	"github.com/Pro7ech/lattigo/mhe"
	"github.com/Pro7ech/lattigo/rlwe"
)

type PrivacyParams struct {
	CKKSParameters hefloat.Parameters `json:"ckks_parameters"`
	Tpk            rlwe.PublicKey     `json:"tpk"`
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

type PublicKeySwitchShare struct {
	Share mhe.KeySwitchingShare `json:"share"`
}
