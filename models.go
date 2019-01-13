package main

type Links struct {
	Href     string
	NoFollow bool
}

var FalseUrls = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
