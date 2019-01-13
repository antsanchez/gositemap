package main

import "encoding/xml"

type Links struct {
	Href     string
	NoFollow bool
}

type UrlSitemap struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

var Extensions = []string{".png", ".jpg", ".jpeg", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db"}

var FalseUrls = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
