package Civets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type civet_uint16 uint16

const (
	RequestMessage uint16 = 1
	ResponseMessage uint16 = 2
)
func (header *civet_uint16) addHeaderProp(value interface{},postion int) *civet_uint16 {
	switch value.(type) {
	case state:
		*header |= (civet_uint16(value.(state))) << postion
	case PathType:
		*header |= (civet_uint16(value.(PathType))) << postion
	case uint8:
		*header |= (civet_uint16(value.(uint8))) << postion
	case uint16:
		*header |= (civet_uint16(value.(uint16))) << postion
	case int16:
		*header |= (civet_uint16(value.(int16))) << postion
	case uint:
		*header |= (civet_uint16(value.(uint))) << postion
	default:
		fmt.Print("\r\nunknow")
	}
	return header
}
func (header civet_uint16) IntToBytes() []byte {
	x := uint16(header)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
func BytesToInt(b []byte) uint {
	bytesBuffer := bytes.NewBuffer(b)
	//fmt.Printf("the b value is %x\n",b)
	//fmt.Printf("the byte value is %x\n",bytesBuffer)
	var x uint16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	//print(x)
	return uint(x)
}
func ProtocolDecode(rec []byte) (*interface{},error){

	reader := bytes.NewReader(rec)
	p := make([]byte, len(rec))
	reader.Read(p)
	fmt.Printf("the rec is %x",rec)
	BufferLength := len(p)
	Head := p[0:2]
	versio := (BytesToInt(Head) & 0x8000) >> 15
	msgType  := uint16((BytesToInt(Head) & 0x6000) >> 13)
	encryption := (BytesToInt(Head) & 0x1000) >> 12
	encryption_type := (BytesToInt(Head) & 0xC00) >> 10
	if versio != 1 {
		return nil,errors.New("Not Civets")
	}

	var temp interface{}
	switch msgType {
	case RequestMessage:
		print("current is request")
		pathType := (BytesToInt(Head) & 0x100) >> 8
		pathLength := BytesToInt(Head) & 0xff
		print("the length is ",pathLength)
		if pathLength <= 0 {
			return nil,errors.New("Path Length Err")
		}
		if int(pathLength) >= BufferLength{
			return nil,errors.New("Path Length Err")
		}
		PathName := p[2:2+pathLength]
		payloadLength := BytesToInt(p[2+pathLength:4+pathLength])
		if int(payloadLength) >= BufferLength{
			return nil,errors.New("Payload Length Err")
		}
		//print("the @ is:",string(p[4+pathLength:5+pathLength]))
		if string(p[4+pathLength:5+pathLength]) != "@" {
			return nil,errors.New("No V")
		}
		payload := p[5+pathLength:5+pathLength+payloadLength]
		if string(p[5+pathLength+payloadLength:6+pathLength+payloadLength]) != "#" {
			return nil,errors.New("Not Civets")
		}
		request := BytesToInt(p[6+pathLength+payloadLength:6+pathLength+payloadLength+2])
		fmt.Printf("\nrequestid is %d\n", request) //68656c6c6f20776f726c64
		timeOut := BytesToInt(p[6+pathLength+payloadLength+2:6+pathLength+payloadLength+4])
		fmt.Printf("\ntime_out is %d\n", timeOut) //68656c6c6f20776f726c64
		temp = CivetRequestPacket{
			uint8(versio),
			uint8(encryption),
			uint8(encryption_type),
			int16(request) ,
			PathType(pathType) ,
			int16(pathLength ),
			PathName,
			int16(payloadLength),
			payload ,
			int16(timeOut)}
	case ResponseMessage:
		temp = CivetResponsePacket{}
		Code := (BytesToInt(Head) & 0xFF)
		fmt.Printf("\ncode is %d\n", Code) //68656c6c6f20776f726c64
		payloadLength := BytesToInt(p[2:4])
		if int(payloadLength) >= BufferLength{
			return nil,errors.New("Not Civets")
		}
		payload := p[5:5+payloadLength]
		fmt.Printf("\npayload is %s\n", payload) //68656c6c6f20776f726c64
		requestId := BytesToInt(p[6+payloadLength:6+payloadLength+2])
		fmt.Printf("\nrequestid is %d\n", requestId) //68656c6c6f20776f726c64
		temp = CivetResponsePacket{
			uint8(versio),
			uint8(encryption),
			uint8(encryption_type) ,
			int16(requestId),
			state(Code),
			int16(payloadLength),
			payload ,
			}
	default:
		print("hi")
	}
	//temp.payloadLength = BytesToInt(p[2:4])
	//temp.Payload = p[5:BytesToInt(c)]
	return &temp,nil
}
func ProtocolEncode(packet *interface{}) (Stream, error) {
	var header civet_uint16 = 0
	var buffer = Stream{}
	switch (*packet).(type) {
	case *CivetRequestPacket:
		//fmt.Printf("\r\nthe value head is %x\r\n",header)
		header.addHeaderProp(version, 15).
			addHeaderProp(RequestMessage, 13).
			addHeaderProp((*packet).(*CivetRequestPacket).encryption, 12).
			addHeaderProp((*packet).(*CivetRequestPacket).encryptionType, 10).
			addHeaderProp((*packet).(*CivetRequestPacket).PathType, 8).
			addHeaderProp((*packet).(*CivetRequestPacket).PathLength, 0)
		buffer.BytesCombine(header.IntToBytes()).
			BytesCombine((*packet).(*CivetRequestPacket).PathName).
			BytesCombine(civet_uint16((*packet).(*CivetRequestPacket).payloadLength).IntToBytes()).
			BytesCombine([]byte("@")).
			BytesCombine((*packet).(*CivetRequestPacket).Payload).
			BytesCombine([]byte("#")).
			BytesCombine(civet_uint16((*packet).(*CivetRequestPacket).RequestId).IntToBytes()).
			BytesCombine(civet_uint16((*packet).(*CivetRequestPacket).timeout).IntToBytes())
		return buffer, nil
	case *CivetResponsePacket:
		header.addHeaderProp(version, 15).
			addHeaderProp(ResponseMessage, 13).
			addHeaderProp((*packet).(*CivetResponsePacket).encryption, 12).
			addHeaderProp((*packet).(*CivetResponsePacket).encryptionType, 10).
			addHeaderProp((*packet).(*CivetResponsePacket).Code, 0)
		buffer.BytesCombine(header.IntToBytes()).
			BytesCombine(civet_uint16((*packet).(*CivetResponsePacket).payloadLength).IntToBytes()).
			BytesCombine([]byte("@")).
			BytesCombine((*packet).(*CivetResponsePacket).Payload).
			BytesCombine([]byte("#")).
			BytesCombine(civet_uint16((*packet).(*CivetResponsePacket).iRequestId).IntToBytes())
		return buffer,nil
	default:
		//fmt.Printf("Param #%d is unknown\n")
		return buffer, errors.New("Param #%d is unknown\n")
	}
}
