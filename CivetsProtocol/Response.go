package Civets

type state uint16

const (
	SERVERSUCCESS       state = 0    //服务器端处理成功
	SERVERDECODEERR     state = 1   //服务器端解码异常
	SERVERENCODEERR     state = 2   //服务器端编码异常
	SERVERNOFUNCERR     state = 3   //服务器端没有该函数
	SERVERNOSERVANTERR  state = 4   //服务器端没有该Servant对象
	SERVERRESETGRID     state = 5   //服务器端灰度状态不一致
	SERVERQUEUETIMEOUT  state = 6   //服务器队列超过限制
	ASYNCCALLTIMEOUT    state = 7   //异步调用超时
	INVOKETIMEOUT       state = 8   //调用超时
	SERVEROVERLOAD      state = 9   //服务器端超负载,超过队列长度
	ADAPTERNULL         state = 10  //客户端选路为空，服务不存在或者所有服务down掉了
	INVOKEBYINVALIDESET state = 11  //客户端按set规则调用非法
	CLIENTDECODEERR     state = 12  //客户端解码异常
	SERVERUNKNOWNERR    state = 13  //服务器端位置异常
)
type CivetResponsePacket struct {
	Version        uint8
	encryption     uint8
	encryptionType uint8
	iRequestId     int16
	Code           state
	payloadLength  int16
	Payload        []byte

}

func CreateResponse(RequestID int16,payload []byte,code state,infos...Handler) *CivetResponsePacket {
	var temp interface{}
	temp = &CivetResponsePacket{
		iRequestId:    RequestID,
		Payload:       payload,
		payloadLength: int16(len(payload)),
		Code:          code,
	}

	for _, info := range infos {
		info.parse(&temp)
	}
	return temp.(*CivetResponsePacket)
}