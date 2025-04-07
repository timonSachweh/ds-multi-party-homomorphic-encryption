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
	PublishPublicKey(urls []string)
	CalculateRelinearizationKeys(urls []string) *rlwe.MemEvaluationKeySet
	Aggregate(weights [][]*rlwe.Ciphertext) []*rlwe.Ciphertext
}

type encryptionServiceImpl struct {
	httpClient           httpclient.DataSpaceClientService
	params               ckks.Parameters
	crs                  *sampling.KeyedPRNG
	encoder              *ckks.Encoder
	publicKeyGenProtocol multiparty.PublicKeyGenProtocol

	publicKeyShareCombined multiparty.PublicKeyGenShare
	publicKey              *rlwe.PublicKey
	evk                    *rlwe.MemEvaluationKeySet
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

	return &encryptionServiceImpl{
		httpClient:           httpClient,
		params:               params,
		crs:                  crs,
		encoder:              ckks.NewEncoder(params),
		publicKeyGenProtocol: multipartyPublicKeyGenProtocol,
		publicKey:            rlwe.NewPublicKey(params),
	}
}

func (e *encryptionServiceImpl) Aggregate(weights [][]*rlwe.Ciphertext) []*rlwe.Ciphertext {
	if len(weights) == 0 {
		return nil
	}
	if len(weights) == 1 {
		return weights[0]
	}
	aggregated := weights[0]
	evaluator := ckks.NewEvaluator(e.params, e.evk)

	for i := range aggregated {
		for j := 1; j < len(weights); j++ {
			if err := evaluator.Mul(aggregated[i], weights[j][i], aggregated[i]); err != nil {
				log.Fatal(err)
			}
			if err := evaluator.Relinearize(aggregated[i], aggregated[i]); err != nil {
				log.Fatal(err)
			}
		}
		err := evaluator.Mul(aggregated[i], 1.0/float64(len(weights)), aggregated[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return aggregated
}

func (e *encryptionServiceImpl) GetInformation() entities.PrivacyParams {
	return entities.PrivacyParams{
		CKKSParameters: e.params,
		//Crs:            e.crs,
	}
}

func (e *encryptionServiceImpl) CalculatePublicKey(clients []string) *rlwe.PublicKey {
	e.calculateSharedCgkShare(clients)
	e.publicKey = rlwe.NewPublicKey(e.params)
	crp := e.publicKeyGenProtocol.SampleCRP(e.crs)
	e.publicKeyGenProtocol.GenPublicKey(e.publicKeyShareCombined, crp, e.publicKey)
	return e.publicKey
}

func (e *encryptionServiceImpl) PublishPublicKey(urls []string) {
	publicKey := entities.PublicKeyExchange{
		PublicKey: *e.publicKey,
	}
	for _, url := range urls {
		err := e.httpClient.SendPublicKeyToClient(url, &publicKey)
		if err != nil {
			log.Fatal(err)
		}
	}
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

func (e *encryptionServiceImpl) CalculateRelinearizationKeys(clients []string) *rlwe.MemEvaluationKeySet {
	rkg := multiparty.NewRelinearizationKeyGenProtocol(e.params)
	_, rkgCombined1, rkgCombined2 := rkg.AllocateShare()

	share := entities.RelinearizationKeyShare{
		ShareOne: rkgCombined1,
		ShareTwo: rkgCombined2,
	}

	for _, client := range clients {
		err := e.httpClient.SendPartialRelinearizationKey(client, &share)
		if err != nil {
			return nil
		}
	}
	rlk := rlwe.NewRelinearizationKey(e.params)
	rkg.GenRelinearizationKey(share.ShareOne, share.ShareTwo, rlk)

	e.evk = rlwe.NewMemEvaluationKeySet(rlk)
	log.Println("Relinearization Key Set created")

	return e.evk
}
