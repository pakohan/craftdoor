package lib

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pakohan/craftdoor/config"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/experimental/devices/mfrc522"
	"periph.io/x/periph/experimental/devices/mfrc522/commands"
	"periph.io/x/periph/host"
)

type Subscriber interface {
	Notify([3]string)
}

type Reader struct {
	s              Subscriber
	deviceFile     string
	rstPin, irqPin gpio.PinIO
	lock           *sync.Mutex
	rlock          *sync.Mutex
	rfid           *mfrc522.Dev
	p              spi.PortCloser
}

func NewReader(cfg config.Config, s Subscriber) (*Reader, error) {
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	rstPinReg := gpioreg.ByName(cfg.RSTPin)
	if rstPinReg == nil {
		return nil, fmt.Errorf("reset pin %s can not be found", cfg.RSTPin)
	}

	irqPinReg := gpioreg.ByName(cfg.IRQPin)
	if irqPinReg == nil {
		return nil, fmt.Errorf("IRQ pin %s can not be found", cfg.IRQPin)
	}

	r := &Reader{
		s:          s,
		deviceFile: cfg.Device,
		rstPin:     rstPinReg,
		irqPin:     irqPinReg,
		lock:       &sync.Mutex{},
		rlock:      &sync.Mutex{},
	}

	err = r.initReader()
	if err != nil {
		return nil, err
	}
	go r.runloop()
	log.Printf("initialized reader")
	return r, nil
}

func (r *Reader) initReader() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.rfid != nil {
		err := r.rfid.Halt()
		if err != nil {
			return err
		}
	}

	if r.p != nil {
		err := r.p.Close()
		if err != nil {
			return err
		}
	}

	var err error
	r.p, err = spireg.Open(r.deviceFile)
	if err != nil {
		return err
	}

	r.rfid, err = mfrc522.NewSPI(r.p, r.rstPin, r.irqPin, mfrc522.WithSync())
	if err != nil {
		return err
	}

	r.rfid.SetAntennaGain(5)
	return nil
}

func (r *Reader) runloop() {
	var old [3]string
	for range time.Tick(1 * time.Second) {
		timeout := 10 * time.Second
		if old[0] != "" {
			timeout = 0
		}

		data, err := r.read(timeout)
		if err != nil {
			log.Printf("err reading data: %s", err)
			continue
		}

		if data != old {
			old = data
			r.s.Notify(data)
		}
	}
}

func (r *Reader) read(timeout time.Duration) ([3]string, error) {
	r.rlock.Lock()
	defer r.rlock.Unlock()

	b0, err := r.rfid.ReadCard(timeout, commands.PICC_AUTHENT1B, 0, 0, mfrc522.DefaultKey)
	if err != nil {
		if err.Error() == "mfrc522 lowlevel: IRQ error" {
			return [3]string{}, r.initReader()
		}

		if strings.HasPrefix(err.Error(), "mfrc522 lowlevel: timeout waiting for IRQ edge: ") {
			return [3]string{}, nil
		}

		return [3]string{}, err
	}

	b1, err := r.rfid.ReadCard(timeout, commands.PICC_AUTHENT1B, 0, 1, mfrc522.DefaultKey)
	if err != nil {
		return [3]string{}, err
	}

	b2, err := r.rfid.ReadCard(timeout, commands.PICC_AUTHENT1B, 0, 2, mfrc522.DefaultKey)
	if err != nil {
		return [3]string{}, err
	}

	return [3]string{
		string(b0),
		string(b1),
		string(b2),
	}, nil
}
