package utils

import (
    "crypto/rand"
)

func RandBytes(size int) (bytes []byte, err error) {
    bytes = make([]byte, 32)
    _, err = rand.Read(bytes)
    return
}
