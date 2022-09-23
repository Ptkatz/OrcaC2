package crypto

import "Orca_Server/pkg/go-engine/crypto/cryptonight"

type Crypto struct {
	cn *cryptonight.CryptoNight
}

func NewCrypto(family string) *Crypto {
	cy := &Crypto{}
	if family == "" || family == "cryptonight" {
		cy.cn = cryptonight.NewCryptoNight()
	}
	return cy
}

func (c *Crypto) Sum(data []byte, algo string, height uint64) []byte {
	return c.cn.Sum(data, algo, height)
}

func TestSum(algo string) bool {
	return cryptonight.TestSum(algo)
}

func TestAllSum() bool {
	for _, algo := range Algo() {
		if !TestSum(algo) {
			return false
		}
	}
	return true
}

func Algo() []string {
	return cryptonight.Algo()
}
