package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"ms_gmail/model"

	"github.com/teamgram/proto/mtproto"
)

func ComputeAESKeyIV(authKey string, plaintext []byte) (aesKey []byte, aesIV []byte, msgKey []byte, err error) {
	// Define x value for client-side messages (x = 0)
	x := 0

	// Validate auth_key length
	if len(authKey) < 88 {
		return nil, nil, nil, errors.New("invalid auth_key length")
	}

	// Extract relevant parts of auth_key
	authKeySubstr1 := []byte(authKey[88+x : 88+x+32])
	authKeySubstr2 := []byte(authKey[x : x+36])
	authKeySubstr3 := []byte(authKey[40+x : 40+x+36])

	// Generate random padding (12 to 1024 bytes, divisible by 16)
	paddingLen := RandomInt(12, 1024)    // Generate random number between 12 and 1024
	paddingLen = (paddingLen + 15) &^ 15 // Ensure padding length is divisible by 16

	randomPadding := generateRandomPadding(int(paddingLen))

	// Calculate msg_key_large
	msgKeyLarge := sha256.Sum256(append(append(authKeySubstr1, plaintext...), randomPadding...))

	// Extract msg_key
	msgKey = msgKeyLarge[8 : 8+16]

	// Calculate sha256_a and sha256_b
	sha256a := sha256.Sum256(append(msgKey, authKeySubstr2...))
	sha256b := sha256.Sum256(append(authKeySubstr3, msgKey...))

	// Extract aes_key and aes_iv
	aesKey = append(sha256a[:8], append(sha256b[8:8+16], sha256a[24:]...)...)
	aesIV = append(sha256b[:8], append(sha256a[8:8+16], sha256b[24:]...)...)

	return aesKey, aesIV, msgKey, nil
}

func ComputeAESKeyIV2(authKey, msgKey []byte) (aesKey []byte, aesIV []byte, err error) {
	// From client x = 0
	x := 0
	
	authKeySubstr2 := []byte(authKey[x : x+36])
	authKeySubstr3 := []byte(authKey[40+x : 40+x+36])

	// Calculate sha256_a and sha256_b
	sha256a := sha256.Sum256(append(msgKey, authKeySubstr2...))
	sha256b := sha256.Sum256(append(authKeySubstr3, msgKey...))

	// Extract aes_key and aes_iv
	aesKey = append(sha256a[:8], append(sha256b[8:8+16], sha256a[24:]...)...)
	aesIV = append(sha256b[:8], append(sha256a[8:8+16], sha256b[24:]...)...)

	return aesKey, aesIV, nil
}

func generateRandomPadding(length int) []byte {
	// Generate random padding of a given length
	padding := make([]byte, length)
	rand.Read(padding)
	return padding
}

// Encrypt plaintext using AES-256 with IEG mode
func EnscriptAES_IGE(aesKey, aesIV, plaintext []byte) ([]byte, error) {
	if len(aesKey) != 32 || len(aesIV) != aes.BlockSize*2 {
		return nil, errors.New("EnscriptAES_IGE: invalid key or iv size")
	}
	if len(plaintext) < aes.BlockSize {
		return nil, errors.New("EnscriptAES_IGE: data too small to encrypt")
	}

	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = ZeroPadding(plaintext, aes.BlockSize)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// Split the 32-byte IV into two 16-byte parts
	t := make([]byte, aes.BlockSize)
	x := make([]byte, aes.BlockSize)
	y := make([]byte, aes.BlockSize)
	copy(x, aesIV[:aes.BlockSize])
	copy(y, aesIV[aes.BlockSize:])

	ciphertext := make([]byte, len(plaintext))
	i := 0
	for i < len(plaintext) {
		xor(x, plaintext[i:i+aes.BlockSize])
		block.Encrypt(t, x)
		xor(t, y)
		x, y = t, plaintext[i:i+aes.BlockSize]
		copy(ciphertext[i:], t)
		i += aes.BlockSize
	}

	return ciphertext, nil
}

func DescriptAES_IGE(aesKey, aesIV, ciphertext []byte) ([]byte, error) {
	// Check key, IV, and ciphertext lengths
	if len(aesKey) != 32 || len(aesIV) != aes.BlockSize*2 {
		return nil, errors.New("DescriptAES_IGE: invalid key, iv")
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("DescriptAES_IGE: data too small to decrypt")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("DescriptAES_IGE: data not divisible by block size")
	}

	block, err := aes.NewCipher(aesKey) // Create new cipher using the key
	if err != nil {
		return nil, err
	}

	// Split the 32-byte IV into two 16-byte parts
	t := make([]byte, aes.BlockSize)
	x := make([]byte, aes.BlockSize)
	y := make([]byte, aes.BlockSize)
	copy(x, aesIV[:aes.BlockSize])
	copy(y, aesIV[aes.BlockSize:])

	plaintext := make([]byte, len(ciphertext))

	i := 0
	for i < len(ciphertext) {
		xor(y, ciphertext[i:i+aes.BlockSize])
		block.Decrypt(t, y)
		xor(t, x)
		y, x = t, ciphertext[i:i+aes.BlockSize]
		copy(plaintext[i:], t)
		i += aes.BlockSize
	}

	// Remove zeros padding
	plaintext = bytes.Trim(plaintext, "\x00")

	return plaintext, nil
}

// Zero padding function
func ZeroPadding(data []byte, blockSize int) []byte {
	paddingLen := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{0}, paddingLen)
	paddedData := append(data, padding...)
	return paddedData
}

func xor(dst, src []byte) {
	for i := range dst {
		dst[i] = dst[i] ^ src[i]
	}
}

// Use SerializeToBuffer2 in core mtproto
func SerializeMarshal(payload model.MessageSending) ([]byte, error) {
	x := mtproto.NewEncodeBuf(32 + int(payload.MessageLen))

	x.Long(payload.Salt)
	x.Long(payload.SessionId)
	x.Long(payload.MessageId)
	x.Int(payload.SeqNo)
	x.Int(payload.MessageLen)
	// x.Bytes([]byte(payload.Body))
	x.Bytes(payload.Body)

	return x.GetBuf(), nil
}

func SerializeUnMarshal(dataBuffer []byte) (*model.MessageSending, error) {
	buf := mtproto.NewDecodeBuf(dataBuffer)
	msg := &model.MessageSending{}

	msg.Salt = buf.Long()
	msg.SessionId = buf.Long()
	msg.MessageId = buf.Long()
	msg.SeqNo = buf.Int()
	msg.MessageLen = buf.Int()

	// Read the remaining bytes as the message body
	// msg.Body = string(buf.Bytes(int(msg.MessageLen)))
	// log.Println("=> msg.MessageLen: ", msg.MessageLen)
	msg.Body = buf.Bytes(int(msg.MessageLen))
	return msg, nil
}

func GetAuthKeyId(authKey string) []byte {
	authKeyBuf := []byte(authKey)
	hash := sha1.Sum(authKeyBuf)
	return hash[len(hash)-8:]
}
