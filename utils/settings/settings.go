package settings

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	GlogConfig   *GlogConfigs
	UDPServer    *UDPServer
	HbaseConfig  *HbaseConfig
	RadiusConfig *RadiusConfig
	SchemaHbase  *SchemaHbase
	ThreadWorker *ThreadWorker
}
type GlogConfigs struct {
	LogDir  string
	MaxSize uint64
	V       int
}
type RadiusConfig struct {
	Secret string
}
type UDPServer struct {
	Host string
	Port int
}
type HbaseConfig struct {
	Host string
	Port int
}
type ThreadWorker struct {
	NumberThreadRadius   int
	NumberThreadIdentity int
}
type ColumFamily1MDO struct {
	Name                string               `json:"Name"`
	QualifierNameCF1MDO *QualifierNameCF1MDO `json:"QualifierNameCF1MDO"`
}

type QualifierNameCF1MDO struct {
	Timestamp string `json:"Timestamp"`
	Type      string `json:"Type"`
	Phone     string `json:"Phone"`
	IPPrivate string `json:"IPPrivate"`
}

type QualifierNameCF1Identity struct {
	Timestamp string `json:"Timestamp"`
	Phone     string `json:"Phone"`
}
type ColumFamily1Identity struct {
	Name                     string                    `json:"Name"`
	QualifierNameCF1Identity *QualifierNameCF1Identity `json:"QualifierNameCF1Identity"`
}
type ColumFamily2Identity struct {
	Name                     string                    `json:"Name"`
	QualifierNameCF2Identity *QualifierNameCF2Identity `json:"QualifierNameCF2Identity"`
}

type QualifierNameCF2Identity struct {
	IPPrivate       string `json:"IPPrivate"`
	PortPrivate     string `json:"PortPrivate"`
	IPDestination   string `json:"IPDestination"`
	PortDestination string `json:"PortDestination"`
}

type MDO struct {
	NameTable       string          `json:"NameTable"`
	ColumFamily1MDO ColumFamily1MDO `json:"ColumFamily1MDO"`
	RadiusTTLHour   int             `json:"RadiusTTLHour"`
}

type Identity struct {
	NameTable            string               `json:"NameTable"`
	ColumFamily1Identity ColumFamily1Identity `json:"ColumFamily1Identity"`
	ColumFamily2Identity ColumFamily2Identity `json:"ColumFamily2Identity"`
	SyslogTTLHour        int                  `json:"SyslogTTLHour"`
}

type SchemaHbase struct {
	MDO      MDO      `json:"MDO"`
	Identity Identity `json:"Identity"`
}

var settings Settings = Settings{}

func init() {
	content, err := ioutil.ReadFile("setting.json")
	if err != nil {
		panic(err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		panic(jsonErr)
	}
}
func GetGlogConfig() *GlogConfigs {
	return settings.GlogConfig
}
func GetUDPServer() *UDPServer {
	return settings.UDPServer
}
func GetHbaseConfig() *HbaseConfig {
	return settings.HbaseConfig
}
func GetRadiusConfig() *RadiusConfig {
	return settings.RadiusConfig
}
func GetSchemaHbase() *SchemaHbase {
	return settings.SchemaHbase
}
func GetThreadWorker() *ThreadWorker {
	return settings.ThreadWorker
}
