package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
	"time"
)

type UUID [16]byte

var (
	mu         sync.Mutex
	lastV7time int64
)

const nanoPerMilli = 1000000

func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// New UUID V4
func New() UUID {
	var uuid UUID
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}

// NewString equals New().String()
func NewString() string {
	return New().String()
}

// getV7Time returns the time in milliseconds and nanoseconds / 256.
// The returned (milli << 12 + seq) is guarenteed to be greater than
// (milli << 12 + seq) returned by any previous call to getV7Time.
func getV7Time() (milli, seq int64) {

	mu.Lock()
	defer mu.Unlock()

	nano := time.Now().UnixNano()
	milli = nano / nanoPerMilli
	// Sequence number is between 0 and 3906 (nanoPerMilli>>8)
	seq = (nano - milli*nanoPerMilli) >> 8
	now := milli<<12 + seq
	if now <= lastV7time {
		now = lastV7time + 1
		milli = now >> 12
		seq = now & 0xfff
	}
	lastV7time = now
	return milli, seq
}

// NewV7 return UUID v7 which is a modernized version of UUID,
// designed specifically for environments that need time-ordered unique identifiers
// while avoiding dependencies on hardware (like MAC addresses in UUID v1).
// It’s built around Unix time, offering millisecond precision,
// which makes it more intuitive and compatible with modern systems than older UUID versions like UUID v1 and v6.
func NewV7() UUID {
	uuid := New()
	_ = uuid[15] // bounds check

	t, s := getV7Time()

	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)

	uuid[6] = 0x70 | (0x0F & byte(s>>8))
	uuid[7] = byte(s)
	return uuid
}

// NewV7String equals NewV7().String()
func NewV7String() string {
	return NewV7().String()
}
