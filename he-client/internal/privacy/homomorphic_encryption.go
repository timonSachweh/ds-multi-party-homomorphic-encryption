package privacy

import (
	"log"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type HEService interface {
	Encrypt(data []float64) (*rlwe.Ciphertext, error)
	Decrypt(ciphertext *rlwe.Ciphertext) ([]float64, error)
}

type HEServiceImpl struct {
	params    ckks.Parameters
	encoder   *ckks.Encoder
	encryptor *rlwe.Encryptor
	decryptor *rlwe.Decryptor
	evaluator *ckks.Evaluator
}

func NewHEService() HEService {
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

	log.Println(params.LogMaxSlots())

	kgen := ckks.NewKeyGenerator(params)
	sk := kgen.GenSecretKeyNew()
	pk := kgen.GenPublicKeyNew(sk)
	encoder := ckks.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, pk)
	decryptor := rlwe.NewDecryptor(params, sk)

	relinearizationKey := kgen.GenRelinearizationKeyNew(sk)
	evaluationKeySet := rlwe.NewMemEvaluationKeySet(relinearizationKey)
	evaluator := ckks.NewEvaluator(params, evaluationKeySet)

	return &HEServiceImpl{
		params:    params,
		encoder:   encoder,
		encryptor: encryptor,
		decryptor: decryptor,
		evaluator: evaluator,
	}
}

func (h *HEServiceImpl) Encrypt(data []float64) (*rlwe.Ciphertext, error) {
	pt := ckks.NewPlaintext(h.params, len(data))

	if err := h.encoder.Encode(data, pt); err != nil {
		log.Fatal(err)
	}

	ciphertext, err := h.encryptor.EncryptNew(pt)
	return ciphertext, err
}

func (h *HEServiceImpl) Decrypt(ciphertext *rlwe.Ciphertext) ([]float64, error) {
	plaintext := h.decryptor.DecryptNew(ciphertext)
	decoded := make([]float64, 5)
	if err := h.encoder.Decode(plaintext, decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
