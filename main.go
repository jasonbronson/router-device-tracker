package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	log.Println("teting")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://192.168.1.254/cgi-bin/devices.ha")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	devices := make([]Device, 0)

	device := Device{}
	doc.Find(".table100 tr").Each(func(i int, s *goquery.Selection) {

		title := s.Find("th").Text()
		detail := s.Find("td").Text()
		titleFormatted := removeNewlines(title)
		detailFormatted := removeNewlines(detail)

		if len(title) > 0 {
			//log.Printf("%v %v \n", titleFormatted, detail)

			if title == "IPv4 Address / Name" {
				nameIPv4 := strings.Split(detail, "/")
				device.IP = removeNewlines(nameIPv4[0])
				device.Name = removeNewlines(nameIPv4[1])
			}

			if title == "Name" {
				device.Name = detailFormatted
			}

			if title == "MAC Address" {
				device.MacAddress = detailFormatted
			}
			if title == "Status" {
				device.Status = detailFormatted
			}

			if strings.Contains(title, "Last Activity") {
				//Router format "Sat May 22 08:29:36 2021"
				lastActivity, err := time.Parse("Mon Jan 2 15:04:05 2006", detailFormatted)
				if err != nil {
					log.Println(err)
				}
				device.LastActivity = lastActivity
			}
			if titleFormatted == "Mesh Client" {
				devices = append(devices, device)
				device = Device{}
			}
		}

	})

	for _, d := range devices {
		
		log.Printf("%v %v %v %v\n", d.Name, d.IP, d.LastActivity.Format("01-02-2006"), d.Status)
	}

}

func removeSpaces(value string) string {
	return strings.ReplaceAll(value, " ", "")
}
func removeNewlines(value string) string {
	return strings.ReplaceAll(value, "\n", "")
}

type Devices struct {
	Device []Device
}

type Device struct {
	Name         string
	LastActivity time.Time
	Status       string
	IP           string
	MacAddress   string
}
