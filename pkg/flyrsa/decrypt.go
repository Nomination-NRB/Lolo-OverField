package flyrsa

import (
	"bytes"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	PublicKey
	keySize int
	n       *big.Int
	d       *big.Int
}

func NewPrivateKey(pemData []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("PEM解码失败")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("无效的PEM类型: %s，期望 RSA PRIVATE KEY", block.Type)
	}
	var keyData struct {
		Version int
		N       *big.Int
		E       *big.Int
		D       *big.Int
	}
	_, err := asn1.Unmarshal(block.Bytes, &keyData)
	if err != nil {
		return nil, fmt.Errorf("ASN.1解析失败: %v", err)
	}
	if keyData.Version != 0 {
		return nil, fmt.Errorf("不支持的私钥版本: %d", keyData.Version)
	}
	return &PrivateKey{
		PublicKey: PublicKey{
			keySize: keyData.N.BitLen(),
			n:       keyData.N,
			e:       keyData.E,
		},
		keySize: keyData.N.BitLen(),
		n:       keyData.N,
		d:       keyData.D,
	}, nil
}

func (priv *PrivateKey) SetND(n, d *big.Int) {
	priv.n = n
	priv.d = d
}

func (priv *PrivateKey) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	var result bytes.Buffer
	reader := bytes.NewReader(ciphertext)
	for reader.Len() > 0 {
		var blockLen int32
		if err := binary.Read(reader, binary.BigEndian, &blockLen); err != nil {
			return nil, err
		}
		encryptedBlock := make([]byte, blockLen)
		if _, err := reader.Read(encryptedBlock); err != nil {
			return nil, err
		}
		decryptedBlock, err := priv.decryptBlock(encryptedBlock)
		if err != nil {
			return nil, err
		}
		result.Write(decryptedBlock)
	}
	return result.Bytes(), nil
}

func (priv *PrivateKey) decryptBlock(encrypted []byte) ([]byte, error) {
	c := new(big.Int).SetBytes(encrypted)
	m := new(big.Int).Exp(c, priv.d, priv.n)
	paddedData := m.Bytes()
	return removePadding(paddedData)
}
