package Civets

type HandleFunc func(*interface{})
type Handler interface {
	parse(packet *interface{})
}

const version uint16 = 1
const (
	path       PathType = 0
	rpc		   PathType = 1
)
const (
	AES  = iota

)

func (f HandleFunc) parse(packet *interface{}) {
	f(packet)
}
func WithPath(str []byte)HandleFunc{
	return func(packet *interface{}) {
		switch (*packet).(type) {
		case *CivetRequestPacket:
			(*packet).(*CivetRequestPacket).PathType = path
			(*packet).(*CivetRequestPacket).PathName = str
			(*packet).(*CivetRequestPacket).PathLength = int16(len(str))
			print("len(str)",len(str))
		default:
			print("error")
		}
	}
}

func WithRPC(str []byte)HandleFunc{
	return func(packet *interface{}) {
		switch (*packet).(type) {
		case *CivetRequestPacket:
			(*packet).(*CivetRequestPacket).PathType = rpc
			(*packet).(*CivetRequestPacket).PathName = str
			(*packet).(*CivetRequestPacket).PathLength = int16(len(str))

		case *CivetResponsePacket:
			(*packet).(*CivetRequestPacket).PathType = rpc
			(*packet).(*CivetRequestPacket).PathName = str
			(*packet).(*CivetRequestPacket).PathLength = int16(len(str))
		default:
			print("error")
		}
	}
}
func WithEncryption(types uint8)HandleFunc{
	return func(packet *interface{}) {
		switch (*packet).(type) {
		case *CivetRequestPacket:
			(*packet).(*CivetRequestPacket).encryption = 1
			(*packet).(*CivetRequestPacket).encryptionType = types

		case *CivetResponsePacket:
			(*packet).(*CivetResponsePacket).encryption = 1
			(*packet).(*CivetResponsePacket).encryptionType = types
		default:
			print("error")
		}
	}
}
func WithTimeout(time int16)HandleFunc{
	return func(packet *interface{}) {
		switch (*packet).(type) {
		case *CivetRequestPacket:
			(*packet).(*CivetRequestPacket).timeout = time
		default:
			print("error")
		}
	}
}