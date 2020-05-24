package lib

import "periph.io/x/periph/experimental/devices/mfrc522"

var craftwerk [16]byte

func init() {
	copy(craftwerk[:], "craftwerk")
}

type Subscriber interface {
	Notify([3]string)
}

type Reader interface {
	InitKey(keyID, keySecret [16]byte, oldKey, keyA, keyB mfrc522.Key) error
}

type dummyReader struct{}

func NewDummyReader() Reader {
	return dummyReader{}
}

func (dr dummyReader) InitKey(keyID, keySecret [16]byte, oldKey, keyA, keyB mfrc522.Key) error {
	return nil
}
