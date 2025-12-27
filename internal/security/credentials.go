package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// CredentialManager encrypts and derives per-user authentication material.
type CredentialManager struct {
	rootKey []byte
	keyID   string
}

// DerivedIdentity provides protocol-agnostic identity fields.
type DerivedIdentity struct {
	AccountID string
	Password  string
	UUID      string
	Username  string
	ID        string
	Secret    string
}

// NewCredentialManager validates config and builds a credential manager.
func NewCredentialManager(masterKey string) (*CredentialManager, error) {
	masterKey = strings.TrimSpace(masterKey)
	if masterKey == "" {
		return nil, fmt.Errorf("credentials: master key is required")
	}
	sum := sha256.Sum256([]byte(masterKey))
	keyID := hex.EncodeToString(sum[:4])
	return &CredentialManager{
		rootKey: sum[:],
		keyID:   keyID,
	}, nil
}

// KeyID returns the identifier for the current master key.
func (m *CredentialManager) KeyID() string {
	return m.keyID
}

// GenerateSecret creates a random credential seed.
func (m *CredentialManager) GenerateSecret() ([]byte, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	return secret, nil
}

// EncryptForUser encrypts a user credential seed and returns ciphertext, nonce, and fingerprint.
func (m *CredentialManager) EncryptForUser(userID uint64, secret []byte) (string, string, string, error) {
	if userID == 0 {
		return "", "", "", fmt.Errorf("credentials: user id required")
	}
	if len(secret) == 0 {
		return "", "", "", fmt.Errorf("credentials: secret required")
	}

	key := m.deriveUserKey(userID, "encryption")
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", "", "", err
	}
	ciphertext := gcm.Seal(nil, nonce, secret, userAAD(userID))
	fingerprint := m.fingerprint(userID, secret)

	return base64.RawStdEncoding.EncodeToString(ciphertext),
		base64.RawStdEncoding.EncodeToString(nonce),
		fingerprint,
		nil
}

// DecryptForUser decrypts a credential seed.
func (m *CredentialManager) DecryptForUser(userID uint64, ciphertext string, nonce string) ([]byte, error) {
	if userID == 0 {
		return nil, fmt.Errorf("credentials: user id required")
	}
	if strings.TrimSpace(ciphertext) == "" || strings.TrimSpace(nonce) == "" {
		return nil, fmt.Errorf("credentials: ciphertext required")
	}

	key := m.deriveUserKey(userID, "encryption")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertextBytes, err := base64.RawStdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	nonceBytes, err := base64.RawStdEncoding.DecodeString(nonce)
	if err != nil {
		return nil, err
	}
	if len(nonceBytes) != gcm.NonceSize() {
		return nil, fmt.Errorf("credentials: invalid nonce")
	}

	return gcm.Open(nil, nonceBytes, ciphertextBytes, userAAD(userID))
}

// DeriveIdentity derives account and password for a user credential version.
func (m *CredentialManager) DeriveIdentity(userID uint64, version int, secret []byte) (DerivedIdentity, error) {
	if userID == 0 {
		return DerivedIdentity{}, fmt.Errorf("credentials: user id required")
	}
	if version <= 0 {
		return DerivedIdentity{}, fmt.Errorf("credentials: version required")
	}
	if len(secret) == 0 {
		return DerivedIdentity{}, fmt.Errorf("credentials: secret required")
	}

	accountBytes := deriveMaterial(secret, userID, version, "account", 16)
	accountID := formatUUID(accountBytes)
	passwordBytes := deriveMaterial(secret, userID, version, "password", 32)
	password := hex.EncodeToString(passwordBytes)

	return DerivedIdentity{
		AccountID: accountID,
		Password:  password,
		UUID:      accountID,
		Username:  accountID,
		ID:        accountID,
		Secret:    password,
	}, nil
}

func (m *CredentialManager) deriveUserKey(userID uint64, purpose string) []byte {
	mac := hmac.New(sha256.New, m.rootKey)
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], userID)
	mac.Write(buf[:])
	mac.Write([]byte(purpose))
	return mac.Sum(nil)
}

func (m *CredentialManager) fingerprint(userID uint64, secret []byte) string {
	key := m.deriveUserKey(userID, "fingerprint")
	mac := hmac.New(sha256.New, key)
	mac.Write(secret)
	return hex.EncodeToString(mac.Sum(nil))
}

func userAAD(userID uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], userID)
	return buf[:]
}

func deriveMaterial(secret []byte, userID uint64, version int, purpose string, size int) []byte {
	mac := hmac.New(sha256.New, secret)
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:8], userID)
	binary.LittleEndian.PutUint64(buf[8:], uint64(version))
	mac.Write(buf[:])
	mac.Write([]byte(purpose))
	sum := mac.Sum(nil)
	if size <= len(sum) {
		return sum[:size]
	}
	return sum
}

func formatUUID(input []byte) string {
	if len(input) < 16 {
		padding := make([]byte, 16-len(input))
		input = append(input, padding...)
	}
	uid := append([]byte(nil), input[:16]...)
	uid[6] = (uid[6] & 0x0f) | 0x40
	uid[8] = (uid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uid[0:4],
		uid[4:6],
		uid[6:8],
		uid[8:10],
		uid[10:16],
	)
}
