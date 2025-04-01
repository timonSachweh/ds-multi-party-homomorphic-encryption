package services

import (
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
	"log"
)

type EncryptionService interface {
	CalculatePublicKey(clients []string) *rlwe.PublicKey
	GetInformation() entities.PrivacyParams
}

type encryptionServiceImpl struct {
	httpClient           httpclient.DataSpaceClientService
	params               ckks.Parameters
	crs                  *sampling.KeyedPRNG
	crp                  multiparty.PublicKeyGenCRP
	encoder              *ckks.Encoder
	publicKeyGenProtocol multiparty.PublicKeyGenProtocol

	publicKeyShareCombined multiparty.PublicKeyGenShare
	publicKey              *rlwe.PublicKey
}

func NewEncryptionService(httpClient httpclient.DataSpaceClientService) EncryptionService {
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
		httpClient:           httpClient,
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

func (e *encryptionServiceImpl) CalculatePublicKey(clients []string) *rlwe.PublicKey {
	e.calculateSharedCgkShare(clients)
	e.publicKey = rlwe.NewPublicKey(e.params)
	e.publicKeyGenProtocol.GenPublicKey(e.publicKeyShareCombined, e.crp, e.publicKey)
	return e.publicKey
}

func (e *encryptionServiceImpl) calculateSharedCgkShare(clients []string) *multiparty.PublicKeyGenShare {
	e.publicKeyShareCombined = e.publicKeyGenProtocol.AllocateShare()
	shareExchange := entities.CkgShareExchange{
		Share: e.publicKeyShareCombined,
	}
	var err error
	for _, client := range clients {
		err = e.httpClient.SendPartialPublicCkgShare(client, &shareExchange)
		if err != nil {
			return nil
		}
	}
	log.Println("Share exchange complete")
	return &e.publicKeyShareCombined
}
