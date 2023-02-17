package radius_utils

import (
	"MobileID/models"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"strings"
	"time"
)

func CreateServerRadius(secret string, packetChRadius chan<- models.Radius) radius.PacketServer {
	//conn, err := udp_utils.GetConnectionUDP(setting.GetUDPInfo().IP, setting.GetUDPInfo().Port)
	//if err != nil {
	//	//fmt.Printf("Some error %v", err)
	//	glog.Error(fmt.Printf("Some error %v", err))
	//	panic(err)
	//}

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		var code radius.Code
		code = radius.CodeAccessAccept
		//log.Printf("\nWriting %v to %v", code, r.RemoteAddr)
		status := rfc2866.AcctStatusType_Get(r.Packet).String()
		phonenumber := rfc2865.CallingStationID_GetString(r.Packet)
		ipprivate := rfc2865.FramedIPAddress_Get(r.Packet).String()
		//message := "RadiusMessage"
		//yyyyMMddHHmmSS
		timestamp := time.Now().Format("20060102150405")
		//line := timestamp + "," + status + "," + phonenumber + "," + ipprivate
		radiusOb := models.Radius{
			Timestamp: timestamp,
			Type:      status,
			Phone:     phonenumber,
			IPPrivate: ipprivate,
		}
		//Loc nhung ban tin ip va phone khac null
		if strings.Compare(ipprivate, "<nil>") != 0 && strings.Compare(phonenumber, "") != 0 && rfc2865.FramedIPAddress_Get(r.Packet).IsPrivate() {
			//println(line)
			packetChRadius <- radiusOb

		}
		w.Write(r.Response(code))
	}
	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(secret)),
	}
	return server
}
