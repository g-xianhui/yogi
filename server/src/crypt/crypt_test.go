package crypt

import (
    "testing"
    "io"
    math_rand "math/rand"
    "crypto/rand"
)

func same(src, dst []byte) bool {
    if len(src) != len(dst) {
        return false
    }

    for i := 0; i < len(src); i++ {
        if src[i] != dst[i] {
            return false
        }
    }

    return true
}

func TestAesEncryt(t *testing.T) {
    key := []byte("example key 1234example key 1234")
    plaintext := []byte("exampleplaintext hahaha")

    ciphertext, err := AesEncrypt(plaintext, key)
    if err != nil {
        t.Error("Encrypt failed: %s", err)
    }

    decrypted, err := AesDecrypt(ciphertext, key)
    if err != nil {
        t.Error("Decrypt failed: %s", err)
    }

    if string(decrypted) != string(plaintext) {
        t.Error("decrypted not the same with plaintext")
    }

    if !same(decrypted, plaintext) {
        t.Error("not the same")
    }
}

func TestWrongKey(t *testing.T) {
    key := []byte("example key 1234example key 1234")
    plaintext := []byte("exampleplaintext")

    ciphertext, err := AesEncrypt(plaintext, key)
    if err != nil {
        t.Error("Encrypt failed: %s", err)
    }

    other_key := []byte("example key 1234example key 1233")
    decrypted, err := AesDecrypt(ciphertext, other_key)
    if err != nil {
        t.Log("Decrypt failed, but that's expected: ", err)
    }

    if same(decrypted, plaintext) {
        t.Error("the same with wrong key")
    }
}

func TestAes(t *testing.T) {
    key := make([]byte, 16)
    for i := 0; i < 100; i++ {
        rand_size := math_rand.Int() % 255 + 1
        plaintext := make([]byte, rand_size)
        if _, err := io.ReadFull(rand.Reader, plaintext); err != nil {
            t.Error("plaintext create failed: ", i, err)
        }

        if _, err := io.ReadFull(rand.Reader, key); err != nil {
            t.Error("key create failed: ", i, err)
        }

        ciphertext, err := AesEncrypt(plaintext, key)
        if err != nil {
            t.Error("AesEncrypt failed: ", i, err)
        }

        decrypted, err := AesDecrypt(ciphertext, key)
        if err != nil {
            t.Error("AesDecrypt failed: ", i, err)
        }

        if !same(decrypted, plaintext) {
            t.Error("not the same: ", i)
        }
    }
}
