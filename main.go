package main

import (
	"MobileID/models"
	"MobileID/udp/server"
	"MobileID/utils/hbase_utils"
	"MobileID/utils/radius_utils"
	"MobileID/utils/settings"
	"MobileID/worker/hbase_worker"
	"MobileID/worker/syslog_worker"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
)

func init() {
	//glog
	//create logs folder
	os.Mkdir("./logs", 0777)
	flag.Lookup("stderrthreshold").Value.Set("[INFO|WARN|FATAL]")
	flag.Lookup("logtostderr").Value.Set("false")
	flag.Lookup("alsologtostderr").Value.Set("true")
	flag.Lookup("log_dir").Value.Set("./logs")
	glog.MaxSize = 1024 * 1024 * settings.GetGlogConfig().MaxSize
	flag.Lookup("v").Value.Set(fmt.Sprintf("%d", settings.GetGlogConfig().V))
	flag.Parse()

}
func main() {
	//Create a channel to receive packet syslog UDP
	packetChIdentity := make(chan []byte)
	packetChRadius := make(chan models.Radius)
	hostUDP := settings.GetUDPServer().Host
	portUDP := settings.GetUDPServer().Port
	connectHbase, err := hbase_utils.GetHBaseClient()
	if err != nil {
		glog.Fatal("Error connect to Hbase cluster: ", err)
	}
	numThreadRadius := settings.GetThreadWorker().NumberThreadRadius
	numThreadIdentity := settings.GetThreadWorker().NumberThreadIdentity
	schemaMDOHbase := settings.GetSchemaHbase().MDO
	schemaIdentityHbase := settings.GetSchemaHbase().Identity
	adminClientHbase, err := hbase_utils.GetAdminHBaseClient()
	hbase_utils.CreateTableMDO(adminClientHbase, schemaMDOHbase)
	//Create table save result after mapping radius with syslog
	hbase_utils.CreateTableIdentity(adminClientHbase, schemaIdentityHbase)
	if err != nil {
		glog.Fatal("Error connect to Admin Hbase cluster: ", err)
	}
	glog.Info("Connected to Hbase cluster")
	server_radius := radius_utils.CreateServerRadius(settings.GetRadiusConfig().Secret, packetChRadius)
	for i := 0; i < numThreadRadius; i++ {
		go hbase_worker.ProcessStreamRadius(connectHbase, schemaMDOHbase, packetChRadius)
	}
	for i := 0; i < numThreadIdentity; i++ {
		go syslog_worker.MappingSyslog(connectHbase, schemaIdentityHbase, schemaMDOHbase, packetChIdentity)
	}
	go func() {
		if err := server_radius.ListenAndServe(); err != nil {
			glog.Fatal(err)
		}
		glog.Info("Started server radius on :1813_receiving radius data")
	}()
	go server.CreateServerUDP(hostUDP, portUDP, packetChIdentity)
	//go func() {
	//	for {
	//		packet := <-packetChIdentity
	//		fmt.Println("Received packet:", string(packet))
	//	}
	//}()
	select {}

}
