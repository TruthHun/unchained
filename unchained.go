// Package unchained implements Django password hashers in Go.
package unchained

import (
    "encoding/base64"
    "fmt"
    "crypto/hmac"
    "crypto/sha1"
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
    "strconv"
    "strings"
)

func b64encode(bytes []byte) string {
    return base64.StdEncoding.EncodeToString(bytes)
}

func compareDigest(val1, val2 string) bool {
    return hmac.Equal([]byte(val1), []byte(val2))
}

func hashAlgorithm(encoded string) string {
    return strings.Split(encoded, "$")[0]
}

// Encode raw password using PBKDF2 SHA256 hasher.
func EncodePBKDF2SHA256(password string, salt string, iterations int) string {
    d := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
    hash := b64encode(d)
    return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iterations, salt, hash)
}

// Validate raw password using PBKDF2 SHA256 hasher.
func VerifyPBKDF2SHA256(password string, encoded string) bool {
    s := strings.Split(encoded, "$")

    if len(s) != 4 {
        return false
    }

    algorithm, iterations, salt := s[0], s[1], s[2]

    if (algorithm != "pbkdf2_sha256") {
        return false
    }

    i, err := strconv.Atoi(iterations)

    if (err != nil) {
        return false
    }

    newencoded := EncodePBKDF2SHA256(password, salt, i)
    return compareDigest(newencoded, encoded)
}

// Encode raw password using PBKDF2 SHA1 hasher.
func EncodePBKDF2SHA1(password string, salt string, iterations int) string {
    d := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha1.Size, sha1.New)
    hash := b64encode(d)
    return fmt.Sprintf("pbkdf2_sha1$%d$%s$%s", iterations, salt, hash)
}

// Validate raw password using PBKDF2 SHA1 hasher.
func VerifyPBKDF2SHA1(password string, encoded string) bool {
    s := strings.Split(encoded, "$")

    if len(s) != 4 {
        return false
    }

    algorithm, iterations, salt := s[0], s[1], s[2]

    if (algorithm != "pbkdf2_sha1") {
        return false
    }

    i, err := strconv.Atoi(iterations)

    if (err != nil) {
        return false
    }

    newencoded := EncodePBKDF2SHA1(password, salt, i)
    return compareDigest(newencoded, encoded)
}

// IsPasswordUsable returns true if encoded password is usable.
func IsPasswordUsable(encoded string) bool {
    if (strings.HasPrefix("!", encoded)) {
        return false
    }

    algorithm := hashAlgorithm(encoded)

    switch algorithm {
        case
            "argon2",
            "bcrypt",
            "bcrypt_sha256",
            "crypt",
            "md5",
            "pbkdf2_sha1",
            "pbkdf2_sha256",
            "sha1",
            "unsalted_md5",
            "unsalted_sha1":
            return true
    }

    return false
}

// CheckPassword validate if the raw password matches the encoded digest.
// This is a shortcut that discovers the algorithm used in the encoded digest
// to perform the correct validation.
func CheckPassword(password string, encoded string) (bool, error) {
    if (!IsPasswordUsable(encoded)) {
        return false, nil
    }

    algorithm := hashAlgorithm(encoded)

    switch algorithm {
        case "pbkdf2_sha256":
            return VerifyPBKDF2SHA256(password, encoded), nil
        case "pbkdf2_sha1":
            return VerifyPBKDF2SHA1(password, encoded), nil
        case
            "argon2",
            "bcrypt",
            "bcrypt_sha256",
            "crypt",
            "md5",
            "sha1",
            "unsalted_md5",
            "unsalted_sha1":
            return false, fmt.Errorf("unchained: hasher not implemented %s", algorithm)
    }

    return false, fmt.Errorf("unchained: invaid hasher %s", algorithm)
}
