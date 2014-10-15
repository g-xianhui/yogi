package crypt

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "bytes"
    "io"
)

func PKCS5Padding(src []byte, block_size int) []byte {
    padding := block_size - len(src) % block_size
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}

func PKCS5UnPadding(src []byte, block_size int) ([]byte, error) {
    length := len(src)
    padding := int(src[length-1])
    if padding < 0 || padding > length {
        return nil, errors.New("PKCS5UnPadding error")
    }
    return src[:length - padding], nil
}

func Encrypt(block cipher.Block, plaintext []byte, key []byte) ([]byte, error) {
    block_size := block.BlockSize()
    plaintext = PKCS5Padding(plaintext, block_size)

    ciphertext := make([]byte, block_size + len(plaintext))
    iv := ciphertext[:block_size]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext[block_size:], plaintext)

    return ciphertext, nil
}

func Decrypt(block cipher.Block, ciphertext []byte, key []byte) ([]byte, error) {
    block_size := block.BlockSize()

    iv := ciphertext[:block_size]
    ciphertext = ciphertext[block_size:]

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(ciphertext, ciphertext)

    return PKCS5UnPadding(ciphertext, block_size)
}

func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return Encrypt(block, plaintext, key)
}

func AesDecrypt(ciphertext []byte, key[]byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return Decrypt(block, ciphertext, key)
}
