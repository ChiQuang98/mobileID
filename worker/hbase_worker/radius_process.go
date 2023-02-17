package hbase_worker

import (
	"MobileID/models"
	"MobileID/utils/hbase_utils"
	"MobileID/utils/settings"
	"github.com/golang/glog"
	"github.com/tsuna/gohbase"
)

// Process channel Radius, save it to Hbase later
func ProcessStreamRadius(clientHbase gohbase.Client, schemaMDO settings.MDO, radiusChan <-chan models.Radius) {
	for radius := range radiusChan {
		//Only save phonenumber with ip's type is start
		if radius.Type == "Start" {
			err := hbase_utils.PutRadiusRecordToHbase(clientHbase, schemaMDO, radius)
			if err != nil {
				glog.Error("Error put radius record to Hbase ", err)
			}
		}

	}
}
