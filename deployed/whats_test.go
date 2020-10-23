package deployed

import (
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	services := []Service{Service{Name: "qa", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"}}

	versions := GetVersion(services)

	got := versions[0]
	want := "npm"

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

	want := []Service{Service{Name: "auth", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"}}
	got, err := LoadServices(r)

	if err != nil {
		log.Fatalf("error loading services from reader: %v", err)
	}

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

	want := []Service{
		Service{Name: "auth", Url: "https://api.cdnjs.com/libraries/jquery", Path: "autoupdate.source"},
		Service{Name: "search", Url: "https://search-service/s", Path: "something"},
	}
	got, err := LoadServices(r)

	if err != nil {
		log.Fatalf("error loading services from reader: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		log.Fatalf("got %v, want %v", got, want)
	}
}

