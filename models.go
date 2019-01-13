package main

type Links struct {
	Href     string
	NoFollow bool
}

var Extensions = []string{".png", ".jpg", ".jpeg", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db"}

var FalseUrls = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
