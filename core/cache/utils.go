package cache

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"time"
)

func fnv32(key any) uint32 {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(fmt.Sprintf("%v", key)))
	return hash.Sum32()
}

func expirationTime(ttl time.Duration) int64 {
	if ttl > 0 {
		return time.Now().Add(ttl).UnixNano()
	}
	return 0
}

// compress compresses data using zlib
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()
	return buf.Bytes(), nil
}

// decompress decompresses data using zlib
func decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	return buf.Bytes(), err
}

func EncodeToHex(value any) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	compressed, err := compress(buf.Bytes())
	return hex.EncodeToString(compressed), err
}

// DecodeFromHex 使用 Gob 反序列化为原始类型
func DecodeFromHex[T any](encodedStr string) (value T, err error) {
	hexBytes, err := hex.DecodeString(encodedStr)
	if err != nil {
		return value, err
	}
	decompressed, err := decompress(hexBytes)
	if err != nil {
		return value, err
	}

	buf := bytes.NewBuffer(decompressed)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&value)
	return value, err
}
