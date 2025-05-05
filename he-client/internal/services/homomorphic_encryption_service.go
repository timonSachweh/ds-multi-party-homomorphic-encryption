package services

import (
	"fmt"
	"github.com/Pro7ech/lattigo/he/hefloat"
	"github.com/Pro7ech/lattigo/mhe"
	"github.com/Pro7ech/lattigo/rlwe"
	"github.com/Pro7ech/lattigo/utils/sampling"
	"github.com/timonSachweh/ds-multi-party-homomorphic-encryption/heclient/internal/entities"
	"log"
	"math"
)

type HEService interface {
	Encrypt(data []float32) ([]*rlwe.Ciphertext, error)
	Decrypt(ciphertext []*rlwe.Ciphertext, vectorLength int) ([]float32, error)
	SetParameters(configuration entities.PrivacyParams)
	PartialShareAggregation(share entities.CkgShareExchange) entities.CkgShareExchange
	PartialRelinKeyAggregation(share entities.RelinearizationKeyShare) entities.RelinearizationKeyShare
	PartialPublicKeySwitchGeneration(weights entities.ClientModel)
	PartialPublicKeySwitchAggregation(share mhe.KeySwitchingShare) mhe.KeySwitchingShare
	SetPublicKey(key entities.PublicKeyExchange)
}

type HEServiceImpl struct {
	params    hefloat.Parameters
	encoder   *hefloat.Encoder
	encryptor *rlwe.Encryptor
	decryptor *rlwe.Decryptor
	evaluator *hefloat.Evaluator

	publicKeyGenProtocol          *mhe.PublicKeyProtocol
	relinearizationKeyGenProtocol *mhe.RelinearizationKeyProtocol
	publicKeySwitchProtocol       *mhe.KeySwitchingProtocol[rlwe.PublicKey]

	secretKey  *rlwe.SecretKey
	ckgShare   *mhe.PublicKeyShare
	relinShare *mhe.RelinearizationKeyShare
	pcksShare  *mhe.KeySwitchingShare
	publicKey  rlwe.EncryptionKey
	seed       [32]byte
	tpk        *rlwe.PublicKey
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
	pt := hefloat.NewPlaintext(h.params, h.params.MaxLevel())

	if err := h.encoder.Encode(data, pt); err != nil {
		return nil, err
	}

	ciphertext := hefloat.NewCiphertext(h.params, 1, h.params.MaxLevel())
	err := h.encryptor.Encrypt(pt, ciphertext)
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
	h.seed = sampling.NewSeed()
	h.params = configuration.CKKSParameters
	h.tpk = &configuration.Tpk
	h.encoder = hefloat.NewEncoder(h.params)

	h.publicKeyGenProtocol = mhe.NewPublicKeyProtocol(h.params)
	h.relinearizationKeyGenProtocol = mhe.NewRelinearizationKeyProtocol(h.params)
	h.publicKeySwitchProtocol = mhe.NewKeySwitchingProtocol[rlwe.PublicKey](h.params)

	h.secretKey = rlwe.NewKeyGenerator(h.params).GenSecretKeyNew()
	h.ckgShare = h.publicKeyGenProtocol.Allocate()
	h.relinShare = h.relinearizationKeyGenProtocol.Allocate()
	h.pcksShare = h.publicKeySwitchProtocol.Allocate(h.params.MaxLevel())

	err := h.publicKeyGenProtocol.Gen(h.secretKey, h.seed, h.ckgShare)
	if err != nil {
		return
	}
}

func (h *HEServiceImpl) SetPublicKey(key entities.PublicKeyExchange) {
	h.publicKey = &key.PublicKey
}

func (h *HEServiceImpl) PartialShareAggregation(share entities.CkgShareExchange) entities.CkgShareExchange {
	err := h.publicKeyGenProtocol.Aggregate(&share.Share, h.ckgShare, h.ckgShare)
	if err != nil {
		return entities.CkgShareExchange{}
	}
	log.Println("partial share aggregation done")
	return share
}

func (h *HEServiceImpl) PartialRelinKeyAggregation(share entities.RelinearizationKeyShare) entities.RelinearizationKeyShare {
	err := h.relinearizationKeyGenProtocol.Aggregate(&share.Share, h.relinShare, &share.Share)
	if err != nil {
		return entities.RelinearizationKeyShare{}
	}
	log.Println("partial relinearization key aggregation done")
	return share
}

func (h *HEServiceImpl) PartialPublicKeySwitchGeneration(weights entities.ClientModel) {
	err := h.publicKeySwitchProtocol.Gen(h.secretKey, h.tpk, float64(1<<30), weights.WeightsAsCiphertext()[0], h.pcksShare)
	if err != nil {
		return
	}
}

func (h *HEServiceImpl) PartialPublicKeySwitchAggregation(share mhe.KeySwitchingShare) mhe.KeySwitchingShare {
	err := h.publicKeySwitchProtocol.Aggregate(&share, h.pcksShare, &share)
	if err != nil {
		return share
	}
	log.Println("partial public key switch aggregation done")
	return share
}
