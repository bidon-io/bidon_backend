package auth

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"hash/crc32"
	"time"
)

var (
	base32NoPadding     = base32.StdEncoding.WithPadding(base32.NoPadding)
	defaultAPIKeyConfig = newAPIKeyConfig("BDNM_", 11) // Bidon Master API
)

type apiKeyConfig struct {
	Prefix string

	EntropyBytes  int
	EntropyDigits int

	UUIDDigits int

	ChecksumBytes  int
	ChecksumDigits int

	Length int
}

// newAPIKeyConfig returns a new apiKeyConfig with the given prefix and entropyBytes.
// If using UUIDv7, entropyBytes should be min 11 to add needed entropy to recommended value of 160.
// UUIDv7 has 74 digits of entropy, 11 bytes give another 88.
func newAPIKeyConfig(prefix string, entropyBytes int) apiKeyConfig {
	c := apiKeyConfig{
		Prefix:        prefix,
		EntropyBytes:  entropyBytes,
		ChecksumBytes: 4, // always 4, crc32.ChecksumIEEE returns 4 bytes
	}

	c.EntropyDigits = base32NoPadding.EncodedLen(c.EntropyBytes)
	c.UUIDDigits = base32NoPadding.EncodedLen(uuid.Size)
	c.ChecksumDigits = base32NoPadding.EncodedLen(c.ChecksumBytes)
	c.Length = len(c.Prefix) + c.EntropyDigits + c.UUIDDigits + c.ChecksumDigits

	return c
}

func NewAPIKey(uuid uuid.UUID) (string, error) {
	c := defaultAPIKeyConfig

	data := make([]byte, 0, c.Length)

	n := copy(data[:len(c.Prefix)], c.Prefix)
	data = data[:n]

	entropy := make([]byte, c.EntropyBytes)
	_, err := rand.Read(entropy)
	if err != nil {
		return "", err
	}
	data = base32NoPadding.AppendEncode(data, entropy)

	data = base32NoPadding.AppendEncode(data, uuid.Bytes())

	checksum := make([]byte, c.ChecksumBytes)
	binary.BigEndian.PutUint32(checksum, crc32.ChecksumIEEE(data))
	data = base32NoPadding.AppendEncode(data, checksum)

	return string(data), nil
}

func ParseAPIKey(key string) (uuid.UUID, error) {
	c := defaultAPIKeyConfig

	if len(key) != c.Length {
		return uuid.Nil, fmt.Errorf("invalid API key length")
	}

	if key[:len(c.Prefix)] != c.Prefix {
		return uuid.Nil, fmt.Errorf("invalid API key prefix")
	}

	data := []byte(key)

	decodedChecksum := make([]byte, c.ChecksumDigits)
	_, err := base32NoPadding.Decode(decodedChecksum, data[len(data)-c.ChecksumDigits:])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid checksum")
	}

	checksum := binary.BigEndian.Uint32(decodedChecksum)
	if checksum != crc32.ChecksumIEEE(data[:len(data)-c.ChecksumDigits]) {
		return uuid.Nil, fmt.Errorf("checksum mismatch")
	}

	uuidBytes := make([]byte, uuid.Size)
	_, err = base32NoPadding.Decode(uuidBytes, data[len(c.Prefix)+c.EntropyDigits:len(data)-c.ChecksumDigits])
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.FromBytes(uuidBytes)
}

type APIKey struct {
	ID uuid.UUID
	User
	PreviousAccessedAt time.Time
}

func (k *APIKey) UserID() int64 {
	return k.User.ID
}

func (k *APIKey) IsAdmin() bool {
	return k.User.IsAdmin
}
