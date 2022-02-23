package Civets

type PathType uint8
type encryptionType uint8


type CivetRequestPacket struct {
	version        uint8
	encryption     uint8
	encryptionType uint8
	RequestId      int16
	PathType       PathType
	PathLength     int16
	PathName       []byte
	payloadLength  int16
	Payload        []byte
	timeout        int16
}

func CreateRequest(RequestID int16,payload []byte,infos...Handler) *CivetRequestPacket {
	var temp interface{}
	temp = &CivetRequestPacket{
		RequestId:     RequestID,
		Payload:       payload,
		payloadLength: int16(len(payload)),
		timeout:       100,
	}

	for _, info := range infos {
		info.parse(&temp)
	}
	return temp.(*CivetRequestPacket)
}