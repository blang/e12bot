package main

import (
	"fmt"
	"net/http"
)

func handleSlotlist(w http.ResponseWriter, r *http.Request) {
	body := ""
	url := ""
	output := ""
	if r.Method == "POST" {
		body = r.PostFormValue("body")
		url = r.PostFormValue("url")
		parsed, err := parseSlotlist(url, body)
		if err != nil {
			output = err.Error()
		} else {
			output = parsed
		}
	}
	fmt.Fprintf(w, "<html><head></head><body><h1>Teste Slotlistparsing</h1>"+
		"Output: <br /><pre>%s</pre><br />"+
		"<form action=\"/slotlist\" method=\"POST\">"+
		"Url: <input type=\"text\" name=\"url\" value=\"%s\" /><br>"+
		"Text: <br /><textarea name=\"body\" cols=\"100\" rows=\"30\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form></body></html>", output, url, body)
}

func parseSlotlist(url string, text string) (string, error) {
	if !parsers.Accept(url) {
		return "", fmt.Errorf("No parser accepts this url: %s", url)
	}

	sl := parsers.Parse(text, url)

	if sl == nil {
		return "", fmt.Errorf("Error while parsing url %s", url)
	}
	parsed := EncodeSlotList(sl)
	return parsed, nil
}
