package auth_test

import (
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	"github.com/gofrs/uuid/v5"
	"testing"
)

func TestNewAPIKey(t *testing.T) {
	id, err := uuid.FromString("0124e053-3580-75a6-973e-a56e7af8b91b")
	if err != nil {
		t.Fatalf("uuid.FromString failed: %v", err)
	}

	key, err := auth.NewAPIKey(id)
	if err != nil {
		t.Fatalf("auth.NewAPIKey failed: %v", err)
	}
	t.Logf("key: %s", key)

	parsedID, err := auth.ParseAPIKey(key)
	if err != nil {
		t.Fatalf("auth.ParseAPIKey failed: %v", err)
	}

	if id != parsedID {
		t.Fatalf("id != parsedID: %s != %s", id, parsedID)
	}
}

func TestParseAPIKey(t *testing.T) {
	tests := map[string]struct {
		key     string
		wantErr bool
	}{
		"valid": {
			key:     "BDNM_B3BOWM3E622ME5W2BYAESOAUZVQB22NFZ6UVXHV6FZDMOZUM7VY",
			wantErr: false,
		},
		"invalid checksum": {
			key:     "BDNM_B3BOWM3E622ME5W2BYAESOAUZVQB22NFZ6UVXHV6FZDMOZUM7V1",
			wantErr: true,
		},
		"invalid length": {
			key:     "BDNM_B3BOWM3E622ME5W2BYAESOAUZVQB22NFZ6UVXHV6FZDMOZUM7",
			wantErr: true,
		},
		"invalid prefix": {
			key:     "BDNX_B3BOWM3E622ME5W2BYAESOAUZVQB22NFZ6UVXHV6FZDMOZUM7VY",
			wantErr: true,
		},
		"key manipulation": {
			key:     "BDNM_B3BOWM3E612ME5W2BYAESOAUZVQB22NFZ6UVXHV6FZDMOZUM7VY",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := auth.ParseAPIKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.ParseAPIKey failed: %v", err)
			}
		})
	}
}
