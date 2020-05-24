package lib

import "periph.io/x/periph/experimental/devices/mfrc522"

var craftwerk [16]byte

func init() {
	copy(craftwerk[:], "craftwerk")
}

// Subscriber is a notification multiplexer
type Subscriber interface {
	Notify([3]string)
}

// Reader accesses a RFID reader
type Reader interface {
	InitKey(keyID, keySecret [16]byte, oldKey, keyA, keyB mfrc522.Key) error
}

type dummyReader struct{}

// NewDummyReader returns a no-op reader implementation for local development
func NewDummyReader() Reader {
	return dummyReader{}
}

func (dr dummyReader) InitKey(keyID, keySecret [16]byte, oldKey, keyA, keyB mfrc522.Key) error {
	return nil
}
