package utils

import (
	"crypto/rand"
	"encoding/hex"
)

type TokenGenerator interface {
    GenerateToken() (string, error)
}

type RandomTokenGenerator struct {
    Size int
}

func NewRandomTokenGenerator(size int) *RandomTokenGenerator {
    return &RandomTokenGenerator{Size: size}
}

func (g *RandomTokenGenerator) GenerateToken() (string, error) {
    if g.Size <= 0 {
        g.Size = 32
    }

    data := make([]byte, g.Size)
    _, err := rand.Read(data)
    if err != nil {
        return "", err
    }

    return hex.EncodeToString(data), nil
}
