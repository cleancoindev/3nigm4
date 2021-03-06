//
// 3nigm4 crypto package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 06/03/2016
//

// Package crypto implements all cryptographic functions
// used by the 3nigm4 suite: i mainly wrap Golang std lib
// function and implement specific pre-processing and
// post-processing logics. This is a security related element
// and should be modified with care: any change to this package
// can potentially modify the security of the whole system.
package crypto

// Golang standard functions
import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

// Extended crypto lib
import (
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
	"golang.org/x/crypto/pbkdf2"
)

// AesMode defines a enum type for available
// aes encryption modes
type AesMode int8

// Available AES modes:
const (
	CBC AesMode = 0 + iota // AES CBC mode
)

var (
	kSaltSize           = 8  // AES salt size
	kRequiredMaxKeySize = 32 // Max key size (AES256)
	kHmacSha256Size     = 32 // Default size for hmac with sha256
)

// PKCS5Padding padding function to pad a certain
// blob of data with necessary data to be used in
// AES block cipher.
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS5UnPadding unpad data after AES block
// decrypting.
func PKCS5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length <= 0 {
		return nil, fmt.Errorf("invalid byte blob lenght: expecting > 0 having %d", length)
	}
	unpadding := int(src[length-1])
	delta := length - unpadding
	if delta < 0 {
		return nil, fmt.Errorf("invalid padding delta lenght: expecting >= 0 having %d", delta)
	}
	return src[:delta], nil
}

// GenerateHMAC produce hmac with a message
// and a key.
func GenerateHMAC(message []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// CheckHMAC verify an hmac message with a given key
// and reference message.
func CheckHMAC(message []byte, messageMAC []byte, key []byte) bool {
	expectedMAC := GenerateHMAC(message, key)
	return hmac.Equal(messageMAC, expectedMAC)
}

// DeriveKeyWithPbkdf2 derive a key from a password using
// Pbkdf2 algorithm. A good number of iterations is
// ~ 10000 cycles. The derivated key has the right
// lenght for being used in AES256.
func DeriveKeyWithPbkdf2(password []byte, salt []byte, iter int) []byte {
	return pbkdf2.Key(password, salt, iter, kRequiredMaxKeySize, sha1.New)
}

// XorKeys xor given keys (passed in a slice)
// returning an unique key.
func XorKeys(keys [][]byte, maxlen int) ([]byte, error) {
	// xor passcodesb
	buffeXored := make([]byte, maxlen)
	for counter, key := range keys {
		if len(key) != maxlen {
			return nil, fmt.Errorf("invalid passcodes: argument passcodes are too short, should be min %d byte long", maxlen)
		}
		// copy or xor
		if counter == 0 {
			copy(buffeXored, key)
		} else {
			for i := 0; i < maxlen; i++ {
				buffeXored[i] ^= key[i]
			}
		}
	}
	return buffeXored, nil
}

// AesEncrypt encrypt data with AES256 using a key.
// Salt and IV will be passed in the encrypted message.
func AesEncrypt(key []byte, salt []byte, plaintext []byte, mode AesMode) ([]byte, error) {
	// check input values
	if len(key) < 1 {
		return nil, fmt.Errorf("invalid key argument: should be not null or empty")
	}
	if len(plaintext) < 1 ||
		plaintext == nil {
		return nil, fmt.Errorf("invalid plain text argument: should be not null or empty")
	}

	// pad plain text
	paddedPlaintext := PKCS5Padding(plaintext, aes.BlockSize)
	// create out buffer
	ciphertext := make([]byte, len(paddedPlaintext)+kSaltSize+aes.BlockSize)
	// copy salt
	if len(salt) != kSaltSize {
		return nil, fmt.Errorf("invalid salt size, expecting %d having %d", kSaltSize, len(salt))
	}
	jdx := 0
	for idx := aes.BlockSize; idx < aes.BlockSize+kSaltSize; idx++ {
		ciphertext[idx] = salt[jdx]
		jdx++
	}

	// Should be previously padded
	if len(paddedPlaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid plain text size: should be a multiple of block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// allocate cipher text buffer
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Select cipher mode
	switch {
	case mode == CBC:
		cipherMode := cipher.NewCBCEncrypter(block, iv)
		cipherMode.CryptBlocks(ciphertext[aes.BlockSize+kSaltSize:], paddedPlaintext)
		break
	}

	// composed as iv + salt + data
	return ciphertext, nil
}

// GetSaltFromCipherText extract the salt component from an
// encrypted data blob.
func GetSaltFromCipherText(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize+kSaltSize {
		return nil, fmt.Errorf("ciphertext is too short: having %d expecting > than %d", len(ciphertext), aes.BlockSize+kSaltSize)
	}
	salt := ciphertext[aes.BlockSize : aes.BlockSize+kSaltSize]
	return salt, nil
}

// AesDecrypt decrypt data with AES256 using a key
// Salt and IV are passed in the encrypted message.
func AesDecrypt(key []byte, ciphertext []byte, mode AesMode) ([]byte, error) {
	// check input values
	if len(key) < 1 ||
		len(ciphertext) < 1 ||
		ciphertext == nil {
		return nil, fmt.Errorf("invalid arguments: should be not null or empty")
	}
	// copy ciphertext to avoid modyfing the actual
	// argument passed data.
	copiedChipertext := make([]byte, len(ciphertext))
	copy(copiedChipertext, ciphertext)

	// get packed values
	iv := copiedChipertext[:aes.BlockSize]
	//salt := ciphertext[aes.BlockSize : aes.BlockSize+kSaltSize]
	ciphert := copiedChipertext[aes.BlockSize+kSaltSize:]

	// check ciphertext lenght
	if len(ciphert) < aes.BlockSize {
		return nil, fmt.Errorf("cipher text too short, must be at least longer than block size")
	}
	if len(ciphert)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("chiper text have wrong size, should be a block size multiple")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Select cipher mode
	switch {
	case mode == CBC:
		cipherMode := cipher.NewCBCDecrypter(block, iv)
		cipherMode.CryptBlocks(ciphert, ciphert)
		break
	}
	// unpad data
	unpadded, err := PKCS5UnPadding(ciphert)
	if err != nil {
		return nil, err
	}

	return unpadded, nil
}

// GetKeyByEmail returns a specific key from an email
// address.
func GetKeyByEmail(keyring openpgp.EntityList, email string) *openpgp.Entity {
	for _, entity := range keyring {
		for _, ident := range entity.Identities {
			if ident.UserId.Email == email {
				return entity
			}
		}
	}
	return nil
}

// OpenPgpEncrypt encrypt using pgp and the passed recipients
// list and signer entity.
func OpenPgpEncrypt(data []byte, recipients openpgp.EntityList, signer *openpgp.Entity) ([]byte, error) {
	// encrypt message
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, recipients, signer, nil, nil)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// OpenPgpDecrypt decrypt a message using the argument
// keyring as source to get required keys.
func OpenPgpDecrypt(data []byte, keyring openpgp.EntityList) ([]byte, error) {
	md, err := openpgp.ReadMessage(bytes.NewBuffer(data), keyring, nil, nil)
	if err != nil {
		return nil, err
	}
	plaintext, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// OpenPgpSignMessage creates a signature for a message.
func OpenPgpSignMessage(msg []byte, signer *openpgp.Entity) ([]byte, error) {
	// new signature struct
	sig := new(packet.Signature)
	sig.SigType = packet.SigTypeBinary
	sig.PubKeyAlgo = signer.PrivateKey.PubKeyAlgo
	sig.CreationTime = time.Now()
	sig.IssuerKeyId = &signer.PrivateKey.KeyId
	sig.Hash = crypto.SHA256

	// generate data hash
	hash := sha256.New()
	io.WriteString(hash, string(msg))

	err := sig.Sign(hash, signer.PrivateKey, nil)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = sig.Serialize(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// OpenPgpVerifySignature verify a signature using a public
// PGP key, an error is returned if the signature is not
// verified otherwise returning nil.
func OpenPgpVerifySignature(signature []byte, message []byte, publicKey *openpgp.Entity) error {
	// load signature
	pack, err := packet.Read(bytes.NewBuffer(signature))
	if err != nil {
		return err
	}
	sign, ok := pack.(*packet.Signature)
	if !ok {
		return fmt.Errorf("unexpected signature format")
	}

	// get signature hash
	hash := sign.Hash.New()

	// hash message content
	_, err = hash.Write(message)
	if err != nil {
		return err
	}

	err = publicKey.PrimaryKey.VerifySignature(hash, sign)
	if err != nil {
		return err
	}

	return nil
}

// Iterate on keys decrypting all encrypted
// entities.
func unlockKeyRing(entity *openpgp.Entity, passphrase []byte) error {
	if entity.PrivateKey != nil &&
		entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt(passphrase)
		if err != nil {
			return err
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt(passphrase)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadArmoredKeyRing read keys in an armored keyring
// and returns openpgp entities. If a passphrase is passed
// it will be used to decrypt keys.
func ReadArmoredKeyRing(kr []byte, passphrase []byte) (openpgp.EntityList, error) {
	// Read armored private key into type EntityList
	// An EntityList contains one or more Entities.
	// This assumes there is only one Entity involved
	kring, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(kr))
	if err != nil {
		return nil, err
	}
	if passphrase != nil {
		for _, entity := range kring {
			err := unlockKeyRing(entity, passphrase)
			if err != nil {
				return nil, err
			}
		}
	}
	return kring, nil
}

const (
	kEn1gm4Type    = "EN1GM4 HANDSHAKE"              // message type;
	kEn1gm4Version = "En1gm4 v1.0.0 (GnuPG v1.4.10)" // Message version.
)

// EncodePgpArmored encode a pgp message in armored
// ASCII format.
func EncodePgpArmored(data []byte, blocktype string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	header := map[string]string{
		"Version": kEn1gm4Version,
	}
	w, err := armor.Encode(buf, blocktype, header)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

// DecodePgpArmored decode pgp armored messages from
// ASCII armored format.
func DecodePgpArmored(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	result, err := armor.Decode(buf)
	if err != nil {
		return nil, err
	}
	decoded, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

var (
	config = packet.Config{
		RSABits:                4096,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		DefaultHash:            crypto.SHA256,
	}
)

func hashToHashId(h crypto.Hash) uint8 {
	v, ok := s2k.HashToHashId(h)
	if !ok {
		panic("tried to convert unknown hash")
	}
	return v
}

// NewPgpKeypair creates a pgp keypair and encodes them as
// byte slides. No encryption is introduced at that point.
func NewPgpKeypair(name, comment, email string) ([]byte, []byte, error) {
	entity, err := openpgp.NewEntity(
		name,
		comment,
		email,
		&config,
	)
	if err != nil {
		return nil, nil, err
	}

	// workaround for issue:
	// https://github.com/golang/go/issues/12153
	for _, id := range entity.Identities {
		id.SelfSignature.PreferredHash = []uint8{hashToHashId(config.DefaultHash)}
	}

	var priv bytes.Buffer
	err = entity.SerializePrivate(&priv, &config)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to serialise private key %s", err.Error())
	}
	var pub bytes.Buffer
	err = entity.Serialize(&pub)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to serialise public key %s", err.Error())
	}

	return priv.Bytes(), pub.Bytes(), nil
}
