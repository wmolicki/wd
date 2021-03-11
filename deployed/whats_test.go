package deployed

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestFetchVersions(t *testing.T) {
	services := []Service{Service{Name: "qa", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"}}

	versions := FetchVersions(services)

	got := versions[0]
	want := Result{Name: "qa", Version: "npm"}

	if got != want {
		t.Fatalf("wanted %s, got %s", want, got)
	}

}

const confStr = `---
services: 
  - 
    name: auth
    path: autoupdate.source
    url: https://api.cdnjs.com/libraries/jquery
`

func TestLoadSingleService(t *testing.T) {
	r := strings.NewReader(confStr)
	c, err := LoadConfig(r)
	if err != nil {
		log.Fatalf("error loading config from reader: %v", err)
	}

	want := []Service{Service{Name: "auth", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"}}
	got := LoadServices(c, "qa")

	if !reflect.DeepEqual(got, want) {
		log.Fatalf("got %v, want %v", got, want)
	}
}

const bigConfStr = `---
services: 
  - 
    name: auth
    path: autoupdate.source
    url: https://api.cdnjs.com/libraries/jquery
  - 
    name: search
    path: something
    url: https://search-service/s
`

func TestLoadManyServices(t *testing.T) {
	r := strings.NewReader(bigConfStr)
	c, err := LoadConfig(r)
	if err != nil {
		log.Fatalf("error loading config from reader: %v", err)
	}

	want := []Service{
		Service{Name: "auth", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"},
		Service{Name: "search", Url: "https://search-service/s", Path: "something"},
	}
	got := LoadServices(c, "qa")

	if !reflect.DeepEqual(got, want) {
		log.Fatalf("got %v, want %v", got, want)
	}
}

func TestConfigFileLoad(t *testing.T) {
	r, err := os.Open("../resources/services_conf.yml")

	if err != nil {
		log.Fatalf("could not open conf file: %v", err)
	}

	c, err := LoadConfig(r)

	if err != nil {
		log.Fatalf("could not load conf: %v", err)
	}

	want := "qa"
	got := c.Def.Env

	if want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}
