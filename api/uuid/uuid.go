package uuid

import "github.com/google/uuid"

type generator func() (uuid.UUID, error)

var googleGenerator = uuid.NewUUID
var googleGeneratorV4 = uuid.NewRandom
var currentGenerator generator
var currentGeneratorV4 generator

func NewUUID() (uuid.UUID, error) {
	return currentGenerator()
}

func NewUUIDv4() (uuid.UUID, error) {
	return currentGeneratorV4()
}

func mockGenerator() (uuid.UUID, error) {
	return [16]byte{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}, nil
}

func IsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func Mock() {
	currentGenerator = mockGenerator
	currentGeneratorV4 = mockGenerator
}

func Unmock() {
	currentGenerator = googleGenerator
	currentGeneratorV4 = googleGeneratorV4
}

func init() {
	currentGenerator = googleGenerator
	currentGeneratorV4 = googleGeneratorV4
}
