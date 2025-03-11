package privacy

import (
	"log"
	"math"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type HEService interface {
	Encrypt(data []float32) ([]*rlwe.Ciphertext, error)
	Encrypt32(data []float32) (*rlwe.Ciphertext, error)
	Encrypt64(data []float64) (*rlwe.Ciphertext, error)
	Decrypt(ciphertext []*rlwe.Ciphertext, vectorLength int) ([]float32, error)
	Decrypt32(ciphertext *rlwe.Ciphertext, vectorLength int) ([]float32, error)
	Decrypt64(ciphertext *rlwe.Ciphertext, vectorLength int) ([]float64, error)
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

func (h *HEServiceImpl) Encrypt(data []float32) ([]*rlwe.Ciphertext, error) {
	converted := make([]float64, len(data))
	for i, v := range data {
		converted[i] = float64(v)
	}
	splits := int(math.Ceil(float64(len(data)) / float64(h.params.MaxSlots())))
	ciphertexts := make([]*rlwe.Ciphertext, splits)
	for i := range int(splits) {
		start := i * h.params.MaxSlots()
		end := min((i+1)*h.params.MaxSlots(), len(data))
		ciphertext, err := h.Encrypt64(converted[start:end])
		if err != nil {
			return nil, err
		}
		ciphertexts[i] = ciphertext
	}
	return ciphertexts, nil
}

func (h *HEServiceImpl) Encrypt32(data []float32) (*rlwe.Ciphertext, error) {
	converted := make([]float64, len(data))
	for i, v := range data {
		converted[i] = float64(v)
	}
	return h.Encrypt64(converted)
}

func (h *HEServiceImpl) Encrypt64(data []float64) (*rlwe.Ciphertext, error) {
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
			decryptionVectorLen = decryptionVectorLen % len(ciphertext)
		}
		dec, err := h.Decrypt64(c, decryptionVectorLen)
		if err != nil {
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

func (h *HEServiceImpl) Decrypt32(ciphertext *rlwe.Ciphertext, vectorLength int) ([]float32, error) {
	decoded, err := h.Decrypt64(ciphertext, vectorLength)
	if err != nil {
		return nil, err
	}

	result := make([]float32, vectorLength)
	for i, v := range decoded {
		result[i] = float32(v)
	}
	return result, nil
}

func (h *HEServiceImpl) Decrypt64(ciphertext *rlwe.Ciphertext, vectorLength int) ([]float64, error) {
	plaintext := h.decryptor.DecryptNew(ciphertext)
	decoded := make([]float64, vectorLength)
	if err := h.encoder.Decode(plaintext, decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
