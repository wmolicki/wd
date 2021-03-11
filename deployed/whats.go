package deployed

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
	Env  string `yaml:"env"`
}

type Defaults struct {
	Env string `yaml:"env"`
}

type Conf struct {
	Svc []Service `yaml:"services"`
	Def Defaults  `yaml:"defaults"`
}

func LoadConfig(reader io.Reader) (*Conf, error) {
	decoder := yaml.NewDecoder(reader)

	c := Conf{}
	err := decoder.Decode(&c)

	if err != nil {
		return nil, fmt.Errorf("could not decode yaml: %v", err)
	}

	return &c, nil
}

func LoadServices(config *Conf, env string) []Service {
	result := make([]Service, 0)
	for _, svc := range config.Svc {
		if svc.Env == env {
			result = append(result, svc)
		}
	}
	return result
}

type Result struct {
	Name    string
	Version string
}

func FetchVersions(services []Service) (result []Result) {

	resultC := make(chan Result)
	// errorsC := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for _, svc := range services {
		wg.Add(1)

		go func(svc Service) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.Url, nil)

			if err != nil {
				log.Printf("could not create request: %v", err)
				return
			}
			log.Printf("calling %s", svc.Url)
			response, err := http.DefaultClient.Do(req)

			if err != nil {
				log.Printf("could not fetch version for %s: %v", svc.Name, err)
				return
			}
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			bodyStr := string(body)

			if err != nil {
				log.Printf("could not read response body: %v", err)
				return
			}

			if !gjson.Valid(bodyStr) {
				log.Printf("response body for %s is not a vald json", svc.Url)
				return
			}

			value := gjson.Get(bodyStr, svc.Path)

			if !value.Exists() {
				log.Printf("could not get value at %s from %s", svc.Path, response.Body)
				return
			}
			if value.IsArray() {
				log.Printf("value at %s is not a single value: %s", svc.Path, response.Body)
				return
			}

			resultC <- Result{Name: svc.Name, Version: value.String()}
		}(svc)
	}

	go func() {
		wg.Wait()
		close(resultC)
	}()

	for {
		select {
		case version, ok := <-resultC:
			if !ok {
				return result
			}

			result = append(result, version)
		}
	}
}
