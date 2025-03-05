package security

import (
	"strings"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		params   *HashParams
	}{
		{
			name:     "default parameters",
			password: "password123",
			params:   nil,
		},
		{
			name:     "custom parameters",
			password: "password123",
			params: &HashParams{
				Time:    2,
				Memory:  32 * 1024,
				Threads: 2,
				KeyLen:  16,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GeneratePasswordHash(tt.password, tt.params)
			if err != nil {
				t.Fatalf("failed to generate password hash: %v", err)
			}

			if !strings.HasPrefix(hash, "$argon2id$") {
				t.Errorf("hash does not have the correct prefix: got %v", hash)
			}

			parts := strings.Split(hash, "$")
			if len(parts) != 6 {
				t.Errorf("hash does not have the correct format: got %v", hash)
			}
		})
	}
}

// TestComparePasswordHash는 입력 비밀번호가 저장된 해시와 일치하는지 확인하는 테스트
// 입력: 평문 비밀번호, 저장된 해시 문자열
// 출력: 일치 여부 (bool), 에러
func TestComparePasswordHash(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		hash        string
		shouldMatch bool
		shouldError bool
	}{
		{
			name:        "valid hash",
			password:    "password123",
			hash:        generateValidHash("password123"),
			shouldMatch: true,
			shouldError: false,
		},
		{
			name:        "invalid hash format",
			password:    "password123",
			hash:        "$argon2id$v=19$m=65536,t=3,p=4$MTIzNDU2Nzg5MDEyMzQ1Ng",
			shouldMatch: false,
			shouldError: true,
		},
		{
			name:        "incorrect password",
			password:    "wrongpassword",
			hash:        generateValidHash("password123"),
			shouldMatch: false,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ComparePasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.shouldError {
				t.Fatalf("expected error: %v, got: %v", tt.shouldError, err)
			}
			if match != tt.shouldMatch {
				t.Errorf("expected match: %v, got: %v", tt.shouldMatch, match)
			}
		})
	}
}

func generateValidHash(password string) string {
	hash, _ := GeneratePasswordHash(password, nil)
	return hash
}
