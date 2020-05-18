package baths

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	validResponseFile = "test/stzh_bath_data.xml"
)

func Test_defaultGetBaths(t *testing.T) {
	baths, err := GetBaths()
	if err != nil {
		t.Errorf("GetBaths failed: %s", err)
	}
	if len(baths) != 16 {
		t.Errorf("Expcted 16 baths, but got %d", len(baths))
	}
}

func Test_customApiGetBaths(t *testing.T) {
	api := ZuerichBathApiClient{
		Client: &http.Client{},
	}
	baths, err := api.GetBaths()
	if err != nil {
		t.Errorf("GetBaths failed: %s", err)
	}
	if len(baths) != 16 {
		t.Errorf("Expcted 16 baths, but got %d", len(baths))
	}
}

func Test_parseXml(t *testing.T) {
	fp, err := os.Open(validResponseFile)
	if err != nil {
		t.Errorf("Failed to open test file %s: %v", validResponseFile, err)
	}
	gotResult, err := parseXml(fp)
	if err != nil {
		t.Errorf("parseXml() error = %v", err)
		return
	}
	if len(gotResult.Baths) != 16 {
		t.Errorf("Expected 16 baths but got %d", len(gotResult.Baths))
	}
	expectedBath := &ZuerichBaths{
		Name:                      "Flussbad Au-Höngg",
		WaterTemperature:          20.0,
		PoiID:                     "flb6938",
		DateModified:              zuerichTime{time.Date(2019, 9, 9, 14, 14, 0, 0, time.UTC)},
		OpenClosedTextPlain:       "geschlossen Saisonende",
		URLPage:                   "https://www.stadt-zuerich.ch/content/ssd/de/index/sport/schwimmen/sommerbaeder/flussbad_au_hoengg.html",
		URLAddressAndOpeningHours: "https://www.stadt-zuerich.ch/content/ssd/de/index/sport/schwimmen/sommerbaeder/flussbad_au_hoengg/jcr:content/mainparsys/texttitleimage.html",
	}
	bath := gotResult.Baths[0]
	if bath.Name != expectedBath.Name {
		t.Errorf("Name doesn't match. Expected %s but got %s", expectedBath.Name, bath.Name)
	}
	if bath.WaterTemperature != expectedBath.WaterTemperature {
		t.Errorf("WaterTemperature doesn't match. Expected %f but got %f", expectedBath.WaterTemperature, bath.WaterTemperature)
	}
	if bath.DateModified != expectedBath.DateModified {
		t.Errorf("DateModified doesn't match. Expected %s but got %s", expectedBath.DateModified, bath.DateModified)
	}
	if bath.OpenClosedTextPlain != expectedBath.OpenClosedTextPlain {
		t.Errorf("OpenClosedTextPlain doesn't match. Expected %s but got %s", expectedBath.OpenClosedTextPlain, bath.OpenClosedTextPlain)
	}
	if bath.PoiID != expectedBath.PoiID {
		t.Errorf("PoiID doesn't match. Expected %s but got %s", expectedBath.PoiID, bath.PoiID)
	}
	if bath.URLAddressAndOpeningHours != expectedBath.URLAddressAndOpeningHours {
		t.Errorf("URLAddressAndOpeningHours doesn't match. Expected\n%s\nbut got\n%s", expectedBath.URLAddressAndOpeningHours, bath.URLAddressAndOpeningHours)
	}
	if bath.URLPage != expectedBath.URLPage {
		t.Errorf("URLPage doesn't match. Expected %s but got %s", expectedBath.URLPage, bath.URLPage)
	}
	if !reflect.DeepEqual(bath, expectedBath) {
		t.Errorf("Deep Equal failed.")
	}
}
