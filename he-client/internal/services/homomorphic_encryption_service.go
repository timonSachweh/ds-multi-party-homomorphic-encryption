package services

import (
	"fmt"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
	"log"
	"math"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type HEService interface {
	Encrypt(data []float32) ([]*rlwe.Ciphertext, error)
	Decrypt(ciphertext []*rlwe.Ciphertext, vectorLength int) ([]float32, error)
	SetParameters(configuration entities.PrivacyParams)
	PartialShareAggregation(share entities.CkgShareExchange) entities.CkgShareExchange
	PartialRelinKeyAggregation(share entities.RelinearizationKeyShare) entities.RelinearizationKeyShare
	SetPublicKey(key entities.PublicKeyExchange)
}

type HEServiceImpl struct {
	params                        ckks.Parameters
	encoder                       *ckks.Encoder
	encryptor                     *rlwe.Encryptor
	decryptor                     *rlwe.Decryptor
	evaluator                     *ckks.Evaluator
	secretKey                     *rlwe.SecretKey
	publicKeyGenProtocol          multiparty.PublicKeyGenProtocol
	ckgShare                      multiparty.PublicKeyGenShare
	relinearizationKeyGenProtocol multiparty.RelinearizationKeyGenProtocol
	rlkEphemSk                    *rlwe.SecretKey
	rkgShareOne                   multiparty.RelinearizationKeyGenShare
	rkgShareTwo                   multiparty.RelinearizationKeyGenShare
	crp                           multiparty.PublicKeyGenCRP
	crs                           sampling.PRNG
	keyCrp                        multiparty.PublicKeyGenCRP
	relinCrp                      multiparty.RelinearizationKeyGenCRP
	publicKey                     rlwe.EncryptionKey
}

func NewHEService() HEService {
	return &HEServiceImpl{}
}

func (h *HEServiceImpl) Encrypt(data []float32) ([]*rlwe.Ciphertext, error) {
	converted := make([]float64, len(data))
	for i, v := range data {
		converted[i] = float64(v)
	}
	splits := int(math.Ceil(float64(len(data)) / float64(h.params.MaxSlots())))
	ciphertexts := make([]*rlwe.Ciphertext, splits)
	h.encryptor = rlwe.NewEncryptor(h.params, h.publicKey)

	for i := range splits {
		start := i * h.params.MaxSlots()
		end := min((i+1)*h.params.MaxSlots(), len(data))
		ciphertext, err := h.encrypt64(converted[start:end])
		if err != nil {
			return nil, err
		}
		ciphertexts[i] = ciphertext
	}
	return ciphertexts, nil
}

func (h *HEServiceImpl) encrypt64(data []float64) (*rlwe.Ciphertext, error) {
	pt := ckks.NewPlaintext(h.params, h.params.MaxLevel())

	if err := h.encoder.Encode(data, pt); err != nil {
		log.Fatal(err)
	}

	ciphertext, err := h.encryptor.EncryptNew(pt)
	return ciphertext, err
}

func (h *HEServiceImpl) Decrypt(ciphertext []*rlwe.Ciphertext, vectorLength int) ([]float32, error) {
	decoded := make([]float64, 0)
	for i, c := range ciphertext {
		decryptionVectorLen := h.params.MaxSlots()
		if i == len(ciphertext)-1 {
			decryptionVectorLen = vectorLength % h.params.MaxSlots()
		}
		dec, err := h.decrypt64(c, decryptionVectorLen)
		if err != nil {
			fmt.Println("Error decrypting ciphertext: ", err)
			return nil, err
		}
		decoded = append(decoded, dec...)
	}

	result := make([]float32, len(decoded))
	for i, v := range decoded {
		result[i] = float32(v)
	}
	return result, nil
}

func (h *HEServiceImpl) decrypt64(ciphertext *rlwe.Ciphertext, vectorLength int) ([]float64, error) {
	plaintext := h.decryptor.DecryptNew(ciphertext)
	decoded := make([]float64, vectorLength)
	if err := h.encoder.Decode(plaintext, decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func (h *HEServiceImpl) SetParameters(configuration entities.PrivacyParams) {
	h.params = configuration.CKKSParameters
	crs, err := sampling.NewKeyedPRNG([]byte{'l', 'a', 't', 't', 'i', 'g', 'o'})
	if err != nil {
		panic(err)
	}
	h.crs = crs
	h.encoder = ckks.NewEncoder(h.params)

	h.publicKeyGenProtocol = multiparty.NewPublicKeyGenProtocol(h.params)
	h.secretKey = rlwe.NewKeyGenerator(h.params).GenSecretKeyNew()
	h.ckgShare = h.publicKeyGenProtocol.AllocateShare()
	h.keyCrp = h.publicKeyGenProtocol.SampleCRP(h.crs)
	h.publicKeyGenProtocol.GenShare(h.secretKey, h.keyCrp, &h.ckgShare)

	h.relinearizationKeyGenProtocol = multiparty.NewRelinearizationKeyGenProtocol(h.params)
	h.relinCrp = h.relinearizationKeyGenProtocol.SampleCRP(h.crs)
	h.rlkEphemSk, h.rkgShareOne, h.rkgShareTwo = h.relinearizationKeyGenProtocol.AllocateShare()
	h.relinearizationKeyGenProtocol.GenShareRoundOne(h.secretKey, h.relinCrp, h.rlkEphemSk, &h.rkgShareOne)
}

func (h *HEServiceImpl) SetPublicKey(key entities.PublicKeyExchange) {
	h.publicKey = &key.PublicKey
}

func (h *HEServiceImpl) PartialShareAggregation(share entities.CkgShareExchange) entities.CkgShareExchange {
	h.publicKeyGenProtocol.AggregateShares(h.ckgShare, share.Share, &h.ckgShare)
	log.Println("partial share aggregation done")
	return share
}

func (h *HEServiceImpl) PartialRelinKeyAggregation(share entities.RelinearizationKeyShare) entities.RelinearizationKeyShare {
	h.relinearizationKeyGenProtocol.AggregateShares(h.rkgShareOne, share.ShareOne, &share.ShareOne)
	h.relinearizationKeyGenProtocol.AggregateShares(h.rkgShareTwo, share.ShareTwo, &share.ShareTwo)
	log.Println("partial relinearization key aggregation done")
	return share
}
