package syslog_worker

import (
	"MobileID/models"
	"MobileID/utils/hbase_utils"
	"MobileID/utils/settings"
	"github.com/golang/glog"
	"github.com/tsuna/gohbase"
	"strings"
	"time"
)

func MappingSyslog(clientHbase gohbase.Client, schemaIdentity settings.Identity, schemaMDO settings.MDO, identityChan <-chan []byte) {
	for record := range identityChan {
		//20230218001800,10.74.178.225,40954,42.1.64.19,443
		attributes := strings.Split(string(record), ",")
		currentTime := time.Now().Format("20060102150405")
		//glog.Info("===CurrentTime: ", currentTime)
		//glog.Info("===syslogtime: ", attributes[0])
		identity := models.Identity{
			Timestamp:       currentTime,
			IPPrivate:       attributes[1],
			PortPrivate:     attributes[2],
			IPDestination:   attributes[3],
			PortDestination: attributes[4],
		}
		err, record_radius := hbase_utils.GetRadiusRecordByRowkey(clientHbase, schemaMDO, identity.IPPrivate)
		//if something wrong or there are no record with the IP private that is provied, just contitnue. ignore behind
		if err != nil {
			continue
		}
		identity.Phone = record_radius.Phone
		timeAccessInternetStr := identity.Timestamp
		timeAssignIPPrivateStr := record_radius.Timestamp
		timeAccessInternet, err := time.Parse("20060102150405", timeAccessInternetStr)
		if err != nil {
			glog.Error("Can not parse timestamp time access Internet web, check again, ", timeAccessInternetStr)
			continue
		}
		timeAssignIPPrivate, err := time.Parse("20060102150405", timeAssignIPPrivateStr)
		if err != nil {
			glog.Error("Can not parse timestamp time time Assign IPPrivate, check again, ", timeAccessInternetStr)
			continue
		}
		//If timestamp phonenumber accessing internet web is after timestamp that phonenumber is assigned ipprivate
		if (record_radius.IPPrivate == identity.IPPrivate) &&
			(timeAccessInternet.After(timeAssignIPPrivate) || timeAccessInternet.Equal(timeAccessInternet)) {
			rowKeyWithOutPort := identity.IPDestination + "|" + identity.Phone
			rowKeyWithPort := identity.IPDestination + "|" + identity.Phone + "|" + identity.PortDestination
			err = hbase_utils.PutIdentityResultRecordToHbase(clientHbase, schemaIdentity, identity, rowKeyWithOutPort)
			if err != nil {
				glog.Error("Error put identity record to Hbase ", err)
			}
			err = hbase_utils.PutIdentityResultRecordToHbase(clientHbase, schemaIdentity, identity, rowKeyWithPort)
			if err != nil {
				glog.Error("Error put identity record to Hbase ", err)
			}
		}
	}
}
