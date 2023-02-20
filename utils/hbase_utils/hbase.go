package hbase_utils

import (
	"MobileID/models"
	"MobileID/utils/settings"
	"context"
	"errors"
	"github.com/golang/glog"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"strconv"
	"time"
)

var client gohbase.Client
var Adminclient gohbase.AdminClient

func GetHBaseClient() (gohbase.Client, error) {
	if client == nil {
		// create a new client if one doesn't exist
		address := settings.GetHbaseConfig().Host + ":" + strconv.Itoa(settings.GetHbaseConfig().Port)
		client = gohbase.NewClient(address)
	}

	// return the singleton client
	return client, nil
}
func GetAdminHBaseClient() (gohbase.AdminClient, error) {
	if Adminclient == nil {
		// create a new client if one doesn't exist
		address := settings.GetHbaseConfig().Host + ":" + strconv.Itoa(settings.GetHbaseConfig().Port)
		Adminclient = gohbase.NewAdminClient(address)
	}

	// return the singleton client
	return Adminclient, nil
}
func CreateTableMDO(adminClient gohbase.AdminClient, schemaMDO settings.MDO) {
	//Check if exist table
	createListTableRequest, err := hrpc.NewListTableNames(context.Background())
	//adminClient
	tableNames, _ := adminClient.ListTableNames(createListTableRequest)
	for _, tb := range tableNames {
		//If table is exsited
		if string(tb.Qualifier) == schemaMDO.NameTable {
			glog.Info("Table MDO is existed, continue processing data")
			return
		}
	}
	table := schemaMDO.NameTable
	families := map[string]map[string]string{schemaMDO.ColumFamily1MDO.Name: {}}
	createTableRequest := hrpc.NewCreateTable(context.Background(), []byte(table), families)
	err = adminClient.CreateTable(createTableRequest)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Info("Created table MDO!!!")
}
func CreateTableIdentity(adminClient gohbase.AdminClient, schemaIdentify settings.Identity) {
	//Check if exist table
	createListTableRequest, err := hrpc.NewListTableNames(context.Background())
	//adminClient
	tableNames, _ := adminClient.ListTableNames(createListTableRequest)
	for _, tb := range tableNames {
		//If table is exsited
		if string(tb.Qualifier) == schemaIdentify.NameTable {
			glog.Info("Table Identitfy is existed, continue processing data")
			return
		}
	}
	tableName := schemaIdentify.NameTable
	families := map[string]map[string]string{schemaIdentify.ColumFamily1Identity.Name: {}, schemaIdentify.ColumFamily2Identity.Name: {}}
	createTableRequest := hrpc.NewCreateTable(context.Background(), []byte(tableName), families)
	err = adminClient.CreateTable(createTableRequest)
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Info("Created table Identity!!!")
}
func GetRadiusRecordByRowkey(clientHbase gohbase.Client, mdoSchema settings.MDO, rowKey string) (error, models.Radius) {
	// create a Get request with the row key
	radius := models.Radius{
		Timestamp: "",
		Type:      "",
		Phone:     "",
		IPPrivate: "",
	}
	tableName := mdoSchema.NameTable
	columnFamily := mdoSchema.ColumFamily1MDO.Name
	getReq, err := hrpc.NewGetStr(context.Background(), tableName, rowKey, hrpc.Families(map[string][]string{columnFamily: nil}))
	if err != nil {
		glog.Error("Error get record radius by rowkey: ", rowKey, err)
		return err, radius
	}

	// send the Get request to HBase
	getRsp, err := clientHbase.Get(getReq)
	if err != nil {
		glog.Error("Error get record radius by rowkey: ", rowKey, err)
		return err, radius
	}
	cells := getRsp.Cells
	if len(cells) < 4 {
		return errors.New("Len of result MDO is not greater than 4, len is " + strconv.Itoa(len(cells))), radius
	}
	// process the row data
	//This is the extractly order of hbase return value
	radius.IPPrivate = string(cells[0].Value)
	radius.Phone = string(cells[1].Value)
	radius.Timestamp = string(cells[2].Value)
	radius.Type = string(cells[3].Value)
	//for _, cell := range getRsp.Cells {
	//
	//	fmt.Printf("column family: %s, column: %s, value: %s\n", cell.Family, cell.Qualifier, cell.Value)
	//}
	return nil, radius
}
func PutRadiusRecordToHbase(clientHbase gohbase.Client, mdoSchema settings.MDO, radius models.Radius) error {
	// Create the Put request
	CF1 := mdoSchema.ColumFamily1MDO
	TTL := time.Duration(mdoSchema.RadiusTTLHour) * time.Hour
	rowKey := radius.IPPrivate
	glog.Info("=RowKey table MDO=====>", rowKey)
	values := map[string]map[string][]byte{CF1.Name: map[string][]byte{
		CF1.QualifierNameCF1MDO.Timestamp: []byte(radius.Timestamp),
		CF1.QualifierNameCF1MDO.Type:      []byte(radius.Type),
		CF1.QualifierNameCF1MDO.Phone:     []byte(radius.Phone),
		CF1.QualifierNameCF1MDO.IPPrivate: []byte(radius.IPPrivate),
	}}
	putRequest, err := hrpc.NewPutStr(context.Background(), mdoSchema.NameTable, rowKey, values, hrpc.TTL(TTL))
	if err != nil {
		glog.Error("Error put record radius to Hbase ", err)
		return err
	}
	_, err = clientHbase.Put(putRequest)
	if err != nil {
		glog.Error("Error put record radius to Hbase ", err)
	}
	//println(rsp.String())
	return nil
}
func PutIdentityResultRecordToHbase(clientHbase gohbase.Client, identitySchema settings.Identity, identity models.Identity) error {
	// Create the Put request
	CF1_iaccess := identitySchema.ColumFamily1Identity
	CF2_idetail := identitySchema.ColumFamily2Identity
	TTL := time.Duration(identitySchema.SyslogTTLHour) * time.Hour
	//Query by ipdestination and phone
	rowKey := identity.IPDestination + "|" + identity.Phone
	glog.Info("=RowKey table identity=====>", rowKey)
	values := map[string]map[string][]byte{
		CF1_iaccess.Name: map[string][]byte{
			CF1_iaccess.QualifierNameCF1Identity.Timestamp: []byte(identity.Timestamp),
			CF1_iaccess.QualifierNameCF1Identity.Phone:     []byte(identity.Phone),
		},
		CF2_idetail.Name: map[string][]byte{
			CF2_idetail.QualifierNameCF2Identity.IPPrivate:       []byte(identity.IPPrivate),
			CF2_idetail.QualifierNameCF2Identity.PortPrivate:     []byte(identity.PortPrivate),
			CF2_idetail.QualifierNameCF2Identity.IPDestination:   []byte(identity.IPDestination),
			CF2_idetail.QualifierNameCF2Identity.PortDestination: []byte(identity.PortDestination),
		},
	}
	putRequest, err := hrpc.NewPutStr(context.Background(), identitySchema.NameTable, rowKey, values, hrpc.TTL(TTL))
	if err != nil {
		glog.Error("Error put record identity to Hbase ", err)
		return err
	}
	_, err = clientHbase.Put(putRequest)
	if err != nil {
		glog.Error("Error put record identity to Hbase ", err)
	}
	//println(rsp.String())
	return nil
}
