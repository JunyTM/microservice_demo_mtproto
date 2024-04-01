package utils

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"ms_auth/infrastructure"
)

const (
	authKeyLen    = 256   // Độ dài key 256 bytes (2048 bits)
	saltLen       = 8     // Độ dài salt 8 bytes (64 bits)
	sessionIdLen  = 8     // Độ dài session ID 8 bytes (64 bits)
	minPaddingLen = 12    // Độ dài padding tối thiểu 12 bytes
	maxPaddingLen = 1024  // Độ dài padding tối đa 1024 bytes
	aesKeyLen     = 32    // Độ dài aes key 32 bytes (256 bits)
	aesIegIvLen   = 32    // Độ dài aes_ieg_iv 32 bytes (256 bits)
	iterations    = 10000 // Số lần lặp PBKDF2
)

func KDF_SHA1(msg_key []byte, salt, sessionId uint64) (aesKey, aesIegIv []byte, err error) {
	authKey := []byte(infrastructure.GetAuthKey())

	// Chuẩn bị dữ liệu cho PBKDF2
	key := append(authKey[:], byte(salt>>56), byte(salt>>48), byte(salt>>40), byte(salt>>32), byte(salt>>24), byte(salt>>16), byte(salt>>8), byte(salt))
	dkLen := aesKeyLen + aesIegIvLen
	prf := sha256.New

	// Tính toán aes_key và aes_ieg_iv
	err = pbkdf2.Key(aesKey[:], aesIegIv[:], prf, key, iterations, dkLen, nil)
	if err != nil {
		return nil, nil, err
	}

	return aesKey, aesIegIv, nil
}

func AES_IEG_Encrypt(ase_key, msg_key string, salt, sessionId []byte, payload interface{}, padding []byte) (data_enscrypted []byte, err error) {
	return nil, nil
}

func CreateMessageKey(authKey []byte, salt, sessionId uint64, payload, padding []byte) []byte {
	// Nối payload và padding
	content := append(payload, padding...)

	// Tạo buffer cho message key
	msgKey := make([]byte, authKeyLen/2+sha256.Size/2)

	// Lấy 32 bytes đầu của auth key
	copy(msgKey[:authKeyLen/2], authKey[:authKeyLen/2])

	// Tính toán SHA-256 của content
	hash := sha256.New()
	hash.Write(content)
	hashSum := hash.Sum(nil)

	// Lấy 128 bit ở giữa của hash (sao chép 64 bytes)
	copy(msgKey[authKeyLen/2:], hashSum[sha256.Size/4:sha256.Size*3/4])

	return msgKey
}
