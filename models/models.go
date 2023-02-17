package models

type Radius struct {
	Timestamp string
	Type      string
	Phone     string
	IPPrivate string
}
type Identity struct {
	Timestamp       string
	Phone           string
	IPPrivate       string
	PortPrivate     string
	IPDestination   string
	PortDestination string
}
