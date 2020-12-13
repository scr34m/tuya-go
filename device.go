package tuya

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Device interface {
	Connect() error
	Get()
}

type device struct {
	ip       string
	id       string
	key      []byte
	c        net.Conn
	sequence int
	ver      string
}

type DevicePayload struct {
	GwID  string      `json:"gwId"`
	DevID string      `json:"devId"`
	T     string      `json:"t"`
	Dps   interface{} `json:"dps"`
	UID   string      `json:"uid"`
}

func NewDevice(id string, key string, ip string) *device {
	return &device{ip: ip, id: id, key: []byte(key), sequence: 0, ver: "3.3"}
}

func (d *device) Get() (interface{}, error) {
	dp := DevicePayload{
		GwID:  d.id,
		DevID: d.id,
		T:     fmt.Sprintf("%d", time.Now().Unix()),
		UID:   d.id,
	}
	x := map[string]string{}
	dp.Dps = x

	dpb, err := json.Marshal(dp)
	if err != nil {
		return nil, err
	}

	pkt, err := d.generateMessage(dpb, CommandTypeDP_QUERY)
	if err != nil {
		return nil, err
	}
	d.send(pkt)

	pkt, err = d.read()
	if err != nil {
		return nil, err
	}

	m := NewMessage(d.key)
	udpm, err := m.Parse(pkt.GetBuffer())
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", udpm)
	return nil, nil
}

func (d *device) generateMessage(data []byte, command int) (Packet, error) {
	payload, err := aesEncrypt(data, d.key)
	if err != nil {
		//return fmt.Errorf("Encrypt error(%v)", err)
		panic(err)
	}

	// TODO !CommandType.DP_QUERY add extra headers to payload
	// TODO work with version != 3.3

	pkt := NewPacket()
	pkt.WriteUInt32BE(uint32(0x000055AA))
	d.sequence++
	pkt.WriteUInt32BE(uint32(d.sequence))
	pkt.WriteUInt32BE(uint32(command))
	pkt.WriteUInt32BE(uint32(len(payload) + (2 * 4)))
	pkt.WriteByte(payload)
	crc := Crc32(pkt.ReadByteRaw(0, len(payload)+(4*4))) & 0xFFFFFFFF
	pkt.WriteUInt32BE(uint32(crc))
	pkt.WriteUInt32BE(uint32(0x0000AA55))
	return pkt, nil
}

func (d *device) Connect() error {
	var err error
	var addr string
	addr = d.ip + ":6668"
	d.c, err = net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		log.Printf("Connection to <%v> failed %v\n", addr, err)
		return err
	}
	return nil
}

func (d *device) read() (Packet, error) {
	header := make([]byte, 4*4)
	_, err := io.ReadFull(d.c, header)
	if err != nil {
		return nil, err
	}

	p := NewPacket()
	p.WriteByte(header)
	p.Rewind()

	if p.ReadUint32BE() != 0x55aa {
		return nil, errors.New("Wrong start mark")
	}

	p.ReadUint32BE() // sequenceN
	p.ReadUint32BE() // commandByte

	payloadSize := int(p.ReadUint32BE())

	response := make([]byte, payloadSize)

	_, err = io.ReadFull(d.c, response)
	if err != nil {
		return nil, err
	}

	p.WriteByte(response)

	return p, nil
}

func (d *device) send(packet Packet) error {
	data := packet.GetBuffer()
	_, err := d.c.Write(data)
	if err != nil {
		return err
	}
	return nil
}
