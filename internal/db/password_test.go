package db

import "testing"

func TestHashPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if hashedPassword == "" {
		t.Fatalf("Expected hashed password, but got empty string")
	}
}

func TestComparePassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, _ := HashPassword(password)

	if result, _ := ComparePassword(hashedPassword, password); !result {
		t.Fatalf("Expected passwords to match, but they didn't")
	}

	if result, _ := ComparePassword(hashedPassword, "wrongPassword"); result {
		t.Fatalf("Expected passwords not to match, but they did")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt, err := generateSalt()

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if len(salt) != 16 {
		t.Fatalf("Expected salt length to be 16, but got: %d", len(salt))
	}
}
