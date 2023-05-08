package uuid

import "github.com/google/uuid"

type generator func() (uuid.UUID, error)

var googleGenerator = uuid.NewUUID
var currentGenerator generator

func NewUUID() (uuid.UUID, error) {
	return currentGenerator()
}

func mockGenerator() (uuid.UUID, error) {
	return [16]byte{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}, nil
}

func Mock() {
	currentGenerator = mockGenerator
}

func Unmock() {
	currentGenerator = googleGenerator
}

func init() {
	currentGenerator = googleGenerator
}
