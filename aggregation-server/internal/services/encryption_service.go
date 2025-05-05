package services

import (
	"fmt"
	"github.com/Pro7ech/lattigo/mhe"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/api/httpclient"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/aggregationserver/internal/entities"

	"github.com/Pro7ech/lattigo/he/hefloat"
	"github.com/Pro7ech/lattigo/ring"
	"github.com/Pro7ech/lattigo/rlwe"
	"log"
)

type EncryptionService interface {
	CalculatePublicKey(clients []string) *rlwe.PublicKey
	GetInformation() entities.PrivacyParams
	PublishPublicKey(urls []string)
	CalculateRelinearizationKeys(urls []string) *rlwe.MemEvaluationKeySet
	Aggregate(weights [][]*rlwe.Ciphertext) []*rlwe.Ciphertext
	CalculatePublicKeySwitchShare(urls []string)
	PublicKeySwitch(weights []*rlwe.Ciphertext) []*rlwe.Ciphertext
	Decrypt(weights []*rlwe.Ciphertext, vectorLen int) []float64
}

type encryptionServiceImpl struct {
	httpClient              httpclient.DataSpaceClientService
	params                  hefloat.Parameters
	encoder                 *hefloat.Encoder
	publicKeyGenProtocol    *mhe.PublicKeyProtocol
	publicKeySwitchProtocol *mhe.KeySwitchingProtocol[rlwe.PublicKey]

	targetPublicKey *rlwe.PublicKey
	targetSecretKey *rlwe.SecretKey

	publicKeyShareCombined  *mhe.PublicKeyShare
	publicKey               *rlwe.PublicKey
	evk                     *rlwe.MemEvaluationKeySet
	publicKeySwitchCombined *mhe.KeySwitchingShare
}

func NewEncryptionService(httpClient httpclient.DataSpaceClientService) EncryptionService {
	var err error
	var params hefloat.Parameters
	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            14,                                    // log2(ring degree)
			LogQ:            []int{55, 45, 45, 45, 45, 45, 45, 45}, // log2(primes Q) (ciphertext modulus)
			LogP:            []int{61},                             // log2(primes P) (auxiliary modulus)
			LogDefaultScale: 45,                                    // log2(scale)
			RingType:        ring.ConjugateInvariant,
		}); err != nil {
		log.Fatal(err)
	}

	multipartyPublicKeyGenProtocol := mhe.NewPublicKeyProtocol(params)
	publicKeySwitchProtocol := mhe.NewKeySwitchingProtocol[rlwe.PublicKey](params)

	targetSecretKey, targetPublicKey := rlwe.NewKeyGenerator(params).GenKeyPairNew()

	return &encryptionServiceImpl{
		httpClient:              httpClient,
		params:                  params,
		encoder:                 hefloat.NewEncoder(params),
		publicKeyGenProtocol:    multipartyPublicKeyGenProtocol,
		targetPublicKey:         targetPublicKey,
		targetSecretKey:         targetSecretKey,
		publicKey:               rlwe.NewPublicKey(params),
		publicKeySwitchProtocol: publicKeySwitchProtocol,
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
	evaluator := hefloat.NewEvaluator(e.params, e.evk)

	for i := range aggregated {
		for j := 1; j < len(weights); j++ {
			if err := evaluator.Add(aggregated[i], weights[j][i], aggregated[i]); err != nil {
				log.Fatal(err)
			}
		}
		err := evaluator.Mul(aggregated[i], 1.0/float64(len(weights)), aggregated[i])
		/*if err := evaluator.Relinearize(aggregated[i], aggregated[i]); err != nil {
			log.Fatal(err)
		}*/
		if err != nil {
			log.Fatal(err)
		}
	}

	return aggregated
}

func (e *encryptionServiceImpl) CalculatePublicKeySwitchShare(clients []string) {
	e.publicKeySwitchCombined = e.publicKeySwitchProtocol.Allocate(e.params.MaxLevel())
	var err error
	for _, client := range clients {
		err = e.httpClient.SendPartialPublicKeySwitchAggregation(client, e.publicKeySwitchCombined)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Public Key Switch Share exchange complete")
}

func (e *encryptionServiceImpl) PublicKeySwitch(weights []*rlwe.Ciphertext) []*rlwe.Ciphertext {
	for i := range weights {
		err := e.publicKeySwitchProtocol.Finalize(weights[i], e.publicKeySwitchCombined, weights[i])
		if err != nil {
			return nil
		}
	}
	return weights
}

func (e *encryptionServiceImpl) Decrypt(weights []*rlwe.Ciphertext, vectorLen int) []float64 {
	decryptor := rlwe.NewDecryptor(e.params, e.targetSecretKey)
	decryptedWeights := make([]float64, 0)
	for i := range weights {
		decVecLen := e.params.MaxSlots()
		if i == len(weights)-1 {
			decVecLen = vectorLen % e.params.MaxSlots()
		}
		dec, err := e.decrypt64(decryptor, weights[i], decVecLen)
		if err != nil {
			fmt.Println("Error decrypting ciphertext: ", err)
			return nil
		}
		decryptedWeights = append(decryptedWeights, dec...)
	}
	return decryptedWeights
}

func (h *encryptionServiceImpl) decrypt64(decryptor *rlwe.Decryptor, ciphertext *rlwe.Ciphertext, vectorLength int) ([]float64, error) {
	plaintext := decryptor.DecryptNew(ciphertext)
	decoded := make([]float64, vectorLength)
	if err := h.encoder.Decode(plaintext, decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func (e *encryptionServiceImpl) GetInformation() entities.PrivacyParams {
	return entities.PrivacyParams{
		CKKSParameters: e.params,
		Tpk:            *e.targetPublicKey,
	}
}

func (e *encryptionServiceImpl) CalculatePublicKey(clients []string) *rlwe.PublicKey {
	e.calculateSharedCgkShare(clients)
	e.publicKey = rlwe.NewPublicKey(e.params)
	err := e.publicKeyGenProtocol.Finalize(e.publicKeyShareCombined, e.publicKey)
	if err != nil {
		return nil
	}
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

func (e *encryptionServiceImpl) calculateSharedCgkShare(clients []string) *mhe.PublicKeyShare {
	e.publicKeyShareCombined = e.publicKeyGenProtocol.Allocate()
	shareExchange := entities.CkgShareExchange{
		Share: *e.publicKeyShareCombined,
	}
	var err error
	for _, client := range clients {
		err = e.httpClient.SendPartialPublicCkgShare(client, &shareExchange)
		if err != nil {
			return nil
		}
	}
	log.Println("Share exchange complete")
	return e.publicKeyShareCombined
}

func (e *encryptionServiceImpl) CalculateRelinearizationKeys(clients []string) *rlwe.MemEvaluationKeySet {
	rkg := mhe.NewRelinearizationKeyProtocol(e.params)
	relinShare := rkg.Allocate()

	share := entities.RelinearizationKeyShare{
		Share: *relinShare,
	}

	for _, client := range clients {
		err := e.httpClient.SendPartialRelinearizationKey(client, &share)
		if err != nil {
			return nil
		}
	}
	rlk := rlwe.NewRelinearizationKey(e.params)
	err := rkg.Finalize(&share.Share, rlk)
	if err != nil {
		return nil
	}

	e.evk = rlwe.NewMemEvaluationKeySet(rlk)
	log.Println("Relinearization Key Set created")

	return e.evk
}
