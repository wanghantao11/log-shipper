package log

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wanghantao11/log-shipper/config"
	"log"
	"net/http"
	"os"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) ParseFile(filename string) ([]Log, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	result := []Log{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var row Log
		err := json.Unmarshal(scanner.Bytes(), &row)
		if err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func (s *Service) AddLog(data []Log) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
	}

	resp, err := http.Post(config.Get(config.ReceiverUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
