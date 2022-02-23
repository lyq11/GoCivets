package Civets

import (
	"bytes"
	//"encoding/binary"
)

type Stream []byte
func (pStream *Stream) BytesCombine(sec Stream) *Stream {
	//fmt.Printf("the sec is %x\n",sec)
	s := make([][]byte, 2)
	s[0] = *pStream
	s[1] = sec
	sep := []byte("")
	*pStream = (*pStream)[0:0]
	//fmt.Printf("the is %x\n",bytes.Join(s, sep))

	for _, value := range bytes.Join(s, sep) {
		//fmt.Printf("the index %d is %x\n",index,value)
		*pStream = append(*pStream, value)
	}
	return pStream
}

