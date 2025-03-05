package security

import (
	"crypto/rand"     // 무작위 솔트 생성용
	"crypto/subtle"   // 타이밍 공격 방지 비교 함수
	"encoding/base64" // 바이너리 데이터를 문자열로 변환
	"errors"          // 커스텀 에러 처리
	"fmt"             // 해시 포맷팅용
	"strings"         // 문자열 파싱용

	"golang.org/x/crypto/argon2" // Argon2id 알고리즘 구현
)

// HashParams는 Argon2의 설정 파라미터를 정의
// - Time: 반복 횟수 (보안 강도에 비례)
// - Memory: 메모리 사용량 (KiB 단위, 메모리 하드 특성 강화)
// - Threads: 병렬 스레드 수 (성능 최적화)
// - KeyLen: 출력 해시 길이 (바이트 단위)
type HashParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// 기본 파라미터 설정
// - Time: 3 (OWASP 최소 권장 2 이상 충족)
// - Memory: 64MB (적당한 보안 수준)
// - Threads: 4 (일반적인 멀티코어 CPU에 적합)
// - KeyLen: 32 (표준 해시 길이)
var defaultParams = HashParams{
	Time:    3,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

// GeneratePasswordHash는 비밀번호를 Argon2로 해싱해 표준화된 문자열 반환
// 입력: 평문 비밀번호, 선택적 파라미터 (nil이면 기본값 사용)
// 출력: "$argon2id$v=버전$m=메모리,t=반복,p=스레드$솔트$해시" 형식, 에러
func GeneratePasswordHash(password string, p *HashParams) (string, error) {

	if len(password) == 0 {
		return "", errors.New("empty password not allowed")
	}

	// 파라미터가 없으면 기본값 사용
	if p == nil {
		p = &defaultParams
	}

	if err := validateParams(p); err != nil {
		return "", err
	}

	// 16바이트 솔트 생성 (무작위성으로 재사용 방지)
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Argon2id로 해시 생성
	// - password: 입력 비밀번호
	// - salt: 무작위 솔트
	// - p.Time, p.Memory, p.Threads, p.KeyLen: 설정된 파라미터
	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	// 솔트와 해시를 base64로 인코딩
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// PHC 형식에 가까운 표준 포맷으로 결합
	// 예: "$argon2id$v=19$m=65536,t=3,p=4$base64salt$base64hash"
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Time, p.Threads, encodedSalt, encodedHash), nil
}

func validateParams(p *HashParams) error {
	// 시간 비용 검증 (OWASP 권장 최소값: 2)
	if p.Time < 2 {
		return errors.New("time cost too low (min: 2)")
	}

	// 메모리 비용 검증 (최소 32MB 권장)
	if p.Memory < 32*1024 {
		return errors.New("memory cost too low (min: 32MB)")
	}

	// 스레드 검증 (최소 1, 최대 임의 제한)
	if p.Threads < 1 {
		return errors.New("parallelism must be at least 1")
	}
	if p.Threads > 64 {
		return errors.New("parallelism exceeds maximum recommended value (max: 64)")
	}

	// 키 길이 검증 (최소 16바이트, 보통 32바이트 권장)
	if p.KeyLen < 16 {
		return errors.New("key length too short (min: 16 bytes)")
	}
	if p.KeyLen > 512 {
		return errors.New("key length exceeds maximum reasonable value (max: 512 bytes)")
	}

	// 메모리와 스레드의 비율 검증 (선택적)
	// 메모리가 충분히 크지 않으면 병렬화 이점이 줄어듦
	if p.Memory < uint32(p.Threads)*8*1024 {
		return errors.New("memory cost should be at least 8MB per thread for efficiency")
	}

	return nil
}

// parseHash는 해시 문자열을 파싱하여 구성 요소를 반환
// 입력: 해시 문자열
// 출력: 메모리, 반복, 스레드, 솔트, 해시, 에러
func parseHash(encodedHash string) (uint32, uint32, uint8, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return 0, 0, 0, nil, nil, errors.New("invalid hash format: must be $argon2id$...")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return 0, 0, 0, nil, nil, fmt.Errorf("unsupported argon2 version: %v", err)
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid parameters: %v", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("failed to decode salt: %v", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("failed to decode hash: %v", err)
	}

	return memory, time, threads, salt, hash, nil
}

// ComparePasswordHash는 입력 비밀번호가 저장된 해시와 일치하는지 확인
// 입력: 평문 비밀번호, 저장된 해시 문자열
// 출력: 일치 여부 (bool), 에러
func ComparePasswordHash(password, encodedHash string) (bool, error) {
	memory, time, threads, salt, hash, err := parseHash(encodedHash)
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
