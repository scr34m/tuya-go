package tuya

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
)

type Message interface {
	Parse([]byte) (*UDPMessage, error)
}

type message struct {
	key []byte
}

func NewMessage(key []byte) Message {
	return &message{key: key}
}

func NewMessageUdp() Message {
	m := md5.New()
	m.Write([]byte(UDP_KEY))
	return &message{key: m.Sum(nil)}
}

const (
	CommandTypeUDP                      = 0
	CommandTypeAP_CONFIG                = 1
	CommandTypeACTIVE                   = 2
	CommandTypeBIND                     = 3
	CommandTypeRENAME_GW                = 4
	CommandTypeRENAME_DEVICE            = 5
	CommandTypeUNBIND                   = 6
	CommandTypeCONTROL                  = 7
	CommandTypeSTATUS                   = 8
	CommandTypeHEART_BEAT               = 9
	CommandTypeDP_QUERY                 = 10
	CommandTypeQUERY_WIFI               = 11
	CommandTypeTOKEN_BIND               = 12
	CommandTypeCONTROL_NEW              = 13
	CommandTypeENABLE_WIFI              = 14
	CommandTypeDP_QUERY_NEW             = 16
	CommandTypeSCENE_EXECUTE            = 17
	CommandTypeUDP_NEW                  = 19
	CommandTypeAP_CONFIG_NEW            = 20
	CommandTypeLAN_GW_ACTIVE            = 240
	CommandTypeLAN_SUB_DEV_REQUEST      = 241
	CommandTypeLAN_DELETE_SUB_DEV       = 242
	CommandTypeLAN_REPORT_SUB_DEV       = 243
	CommandTypeLAN_SCENE                = 244
	CommandTypeLAN_PUBLISH_CLOUD_CONFIG = 245
	CommandTypeLAN_PUBLISH_APP_CONFIG   = 246
	CommandTypeLAN_EXPORT_APP_CONFIG    = 247
	CommandTypeLAN_PUBLISH_SCENE_PANEL  = 248
	CommandTypeLAN_REMOVE_GW            = 249
	CommandTypeLAN_CHECK_GW_UPDATE      = 250
	CommandTypeLAN_GW_UPDATE            = 251
	CommandTypeLAN_SET_GW_CHANNEL       = 252
)

var UDP_KEY = "yGAdlopoPVldABfn"

func (m *message) Parse(buffer []byte) (*UDPMessage, error) {
	p := NewPacketFromByte(buffer)

	// check start marker
	if p.ReadUint32BE() != 0x55aa {
		return nil, errors.New("Starting marker wrong")
	}

	p.ReadUint32BE() // sequenceN
	p.ReadUint32BE() // commandByte

	payloadSize := int(p.ReadUint32BE())

	if len(buffer)-8 < payloadSize {
		return nil, errors.New("Packet missing payload")
	}

	// TODO returnCode check
	p.ReadUint32BE()

	// payload with out end mardker
	payload := p.ReadByte(payloadSize - 4 - 4 - 4)

	// TODO CRC check
	crc := p.ReadUint32BE()
	if crc != Crc32(buffer[:payloadSize+8]) {
		return nil, errors.New("CRC error")
	}

	// check end marker
	if p.ReadUint32BE() != 0xaa55 {
		return nil, errors.New("End marker wrong")
	}

	message, err := aesDecrypt(payload, m.key)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}

	var u UDPMessage
	json.Unmarshal(message, &u)

	return &u, nil
}
