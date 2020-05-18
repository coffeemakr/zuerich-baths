package baths

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	freibadZuerichUrl = "https://www.stadt-zuerich.ch/stzh/bathdatadownload"
	userAgent         = "badi/0.1 (+https://github.com/coffeemakr/zuerich-baths)"
)

type ZuerichBaths struct {
	Name                string      `xml:"title"`
	WaterTemperature    float32     `xml:"temperatureWater"`
	PoiID               string      `xml:"poiid"`
	DateModified        zuerichTime `xml:"dateModified"`
	OpenClosedTextPlain string      `xml:"openClosedTextPlain"`
	URLPage             string      `xml:"urlPage"`
	//PathPage                  string `xml:"pathPage"`
	URLAddressAndOpeningHours string `xml:"urlAddressAndOpeningHours"`
}

type zuerichTime struct {
	time.Time
}

func (c *zuerichTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	second := strings.Split(v, ",")[1]
	parse, err := time.Parse(" 02.01.2006 15:04", second)
	if err != nil {
		return err
	}
	*c = zuerichTime{parse}
	return nil
}

func (b ZuerichBaths) String() string {
	return fmt.Sprintf("Bath{Name: \"%s\", WaterTemperature: %f, PoiID: \"%s\", DateModified: \"%s\", OpenClosedTextPlain: \"%s\", URLPage: \"%s\", URLAddressAndOpeningHours: \"%s\"}",
		b.Name,
		b.WaterTemperature,
		b.PoiID,
		b.DateModified,
		b.OpenClosedTextPlain,
		b.URLPage,
		b.URLAddressAndOpeningHours,
	)
}

type rawZuerichBashData struct {
	Baths []*ZuerichBaths `xml:"baths>bath"`
}

func parseXml(r io.Reader) (result *rawZuerichBashData, err error) {
	result = new(rawZuerichBashData)
	err = xml.NewDecoder(r).Decode(result)
	return
}

type ZuerichBathApiClient struct {
	Client *http.Client
}

var defaultApiClient = ZuerichBathApiClient{
	Client: http.DefaultClient,
}

// GetBaths returns a list of public open-air baths in Zürich
func (a ZuerichBathApiClient) GetBaths() ([]*ZuerichBaths, error) {
	req, err := http.NewRequest(http.MethodGet, freibadZuerichUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml")
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	result, err := parseXml(resp.Body)
	if err != nil {
		return nil, err
	}
	return result.Baths, nil
}

// GetBaths returns a list of public open-air baths in Zürich using the default API client
func GetBaths() ([]*ZuerichBaths, error) {
	return defaultApiClient.GetBaths()
}
