package services

import (
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
	"log"
)

type EncryptionService interface {
	GetInformation() entities.PrivacyParams
}

type encryptionServiceImpl struct {
	params               ckks.Parameters
	crs                  *sampling.KeyedPRNG
	crp                  multiparty.PublicKeyGenCRP
	encoder              *ckks.Encoder
	publicKeyGenProtocol multiparty.PublicKeyGenProtocol

	publicKeyCombined multiparty.PublicKeyGenShare
	publicKey         *rlwe.PublicKey
}

func NewEncryptionService() EncryptionService {
	var err error
	var params ckks.Parameters
	if params, err = ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            14,                                    // log2(ring degree)
		LogQ:            []int{55, 45, 45, 45, 45, 45, 45, 45}, // log2(primes Q) (ciphertext modulus)
		LogP:            []int{61},                             // log2(primes P) (auxiliary modulus)
		LogDefaultScale: 45,                                    // log2(scale)
		RingType:        ring.ConjugateInvariant,
	}); err != nil {
		log.Fatal(err)
	}

	crs, err := sampling.NewKeyedPRNG([]byte{'l', 'a', 't', 't', 'i', 'g', 'o'})
	if err != nil {
		panic(err)
	}

	multipartyPublicKeyGenProtocol := multiparty.NewPublicKeyGenProtocol(params)
	crp := multipartyPublicKeyGenProtocol.SampleCRP(crs)

	return &encryptionServiceImpl{
		params:               params,
		crs:                  crs,
		crp:                  crp,
		encoder:              ckks.NewEncoder(params),
		publicKeyGenProtocol: multipartyPublicKeyGenProtocol,
		publicKey:            rlwe.NewPublicKey(params),
	}
}

func (e *encryptionServiceImpl) GetInformation() entities.PrivacyParams {
	return entities.PrivacyParams{
		CKKSParameters:       e.params,
		PublicKeyGenProtocol: e.publicKeyGenProtocol,
		Crp:                  e.crp,
	}
}
