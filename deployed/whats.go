package deployed

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
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
}

type Conf struct {
	Svc []Service `yaml:"services"`
}

func LoadServices(reader io.Reader) ([]Service, error) {
	decoder := yaml.NewDecoder(reader)

	c := Conf{}
	err := decoder.Decode(&c)

	if err != nil {
		return nil, fmt.Errorf("could not decode yaml: %v", err)
	}

	return c.Svc, nil
}

func GetVersion(services []Service) (result []string) {

	resultC := make(chan string)
	// errorsC := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
			time.Sleep(time.Duration(rand.Intn(6)) * time.Second)

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

			resultC <- value.String()
		}(svc)
	}

	// for version := range resultC {
	// 	result = append(result, version)
	// }

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
