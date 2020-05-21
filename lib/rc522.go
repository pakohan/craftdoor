package lib

import (
	"fmt"
	"log"
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
	Notify(string)
}

type Reader struct {
	s              Subscriber
	deviceFile     string
	rstPin, irqPin gpio.PinIO
	lock           *sync.Mutex
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
	}

	return r, r.initReader()
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
	log.Printf("initialized reader")
	return nil
}

func (r *Reader) runloop() {
	var old string
	for range time.Tick(1 * time.Second) {
		data, err := r.rfid.ReadCard(10*time.Second, commands.PICC_AUTHENT1B, 0, 0, mfrc522.DefaultKey)
		if err != nil {
			if err.Error() == "mfrc522 lowlevel: IRQ error" {
				err = r.initReader()
				if err != nil {
					log.Println("err initializing pin after error")
				}
			} else {
				log.Printf("err from ReadCard: %s", err)
			}
		} else if old != string(data) {
			old = string(data)
			r.s.Notify(old)
			log.Printf("state changed to '%s'", old)
		}
	}
}
