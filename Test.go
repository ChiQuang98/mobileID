package main

import (
	"MobileID/utils/hbase_utils"
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/tsuna/gohbase/hrpc"
)

func main() {
	connectHbase, err := hbase_utils.GetHBaseClient()
	if err != nil {
		glog.Fatal("Error connect to Hbase cluster: ", err)
	}
	// specify the table name and row key
	tableName := "identity"
	rowKey := "42.1.64.19|quang"
	columnFamily := "idetail"

	// create a Get request with the row key
	getReq, err := hrpc.NewGetStr(context.Background(), tableName, rowKey, hrpc.Families(map[string][]string{columnFamily: nil}))
	if err != nil {
		fmt.Println(err)
		return
	}

	// send the Get request to HBase
	getRsp, err := connectHbase.Get(getReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	// process the row data
	for _, cell := range getRsp.Cells {
		fmt.Printf("column family: %s, column: %s, value: %s\n", cell.Family, cell.Qualifier, string(cell.Value))
	}

	// close the client connection
	connectHbase.Close()
}
