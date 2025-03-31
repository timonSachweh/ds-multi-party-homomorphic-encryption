package entities

import (
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type PrivacyParams struct {
	CKKSParameters       ckks.Parameters                 `json:"ckks_parameters"`
	PublicKeyGenProtocol multiparty.PublicKeyGenProtocol `json:"public_key_gen_protocol"`
	Crp                  multiparty.PublicKeyGenCRP      `json:"crp"`
}

type CkgShareExchange struct {
	Share multiparty.PublicKeyGenShare `json:"share"`
}
