package flyrsa

import (
	"bytes"
	"crypto/rand"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
)

type PublicKey struct {
	keySize int
	n       *big.Int
	e       *big.Int
}

func NewPublicKey(pemData []byte) (*PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("PEM解码失败")
	}
	if block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("无效的PEM类型: %s，期望 RSA PUBLIC KEY", block.Type)
	}
	var keyData struct {
		N *big.Int
		E *big.Int
	}
	_, err := asn1.Unmarshal(block.Bytes, &keyData)
	if err != nil {
		return nil, fmt.Errorf("ASN.1解析失败: %v", err)
	}
	return &PublicKey{keySize: keyData.N.BitLen(), n: keyData.N, e: keyData.E}, nil
}

func (publ *PublicKey) SetNE(n, e *big.Int) {
	publ.n = n
	publ.e = e
}

func (publ *PublicKey) Encode(data []byte) ([]byte, error) {
	blockSize := publ.n.BitLen() / 8
	maxDataSize := blockSize - 11
	var result bytes.Buffer

	for i := 0; i < len(data); i += maxDataSize {
		end := i + maxDataSize
		if end > len(data) {
			end = len(data)
		}
		blockData := data[i:end]
		encryptedBlock, err := publ.encryptBlock(blockData, blockSize)
		if err != nil {
			return nil, err
		}
		binary.Write(&result, binary.BigEndian, int32(len(encryptedBlock)))
		result.Write(encryptedBlock)
	}

	return result.Bytes(), nil
}

func (publ *PublicKey) encryptBlock(data []byte, blockSize int) ([]byte, error) {
	paddedData := publ.addPadding(data, blockSize)
	m := new(big.Int).SetBytes(paddedData)
	if m.Cmp(publ.n) >= 0 {
		return nil, fmt.Errorf("message must be smaller than the modulus")
	}
	c := new(big.Int).Exp(m, publ.e, publ.n)

	return c.Bytes(), nil
}

func (publ *PublicKey) addPadding(data []byte, blockSize int) []byte {
	if len(data) > blockSize-1 {
		panic("message too large")
	}
	padded := make([]byte, blockSize)
	padded[0] = 0x01
	binary.BigEndian.PutUint32(padded[1:5], uint32(len(data)))
	randomBytes := make([]byte, blockSize-5-len(data))
	rand.Read(randomBytes)
	for i, b := range randomBytes {
		signedByte := int8(b)
		padded[5+i] = byte(signedByte)
	}
	copy(padded[blockSize-len(data):], data)
	return padded
}
