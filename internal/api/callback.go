package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type load struct {
	Token    string
	Decision string
	Comment  string
}

func SendCallback(callbackurl, token, decision, comment string) error {

	var l load
	l.Comment = comment
	l.Decision = decision
	l.Token = token

	var j, err = json.Marshal(l)
	if err != nil {
		log.Printf("Error making json format: %v", err)
		return err
	}

	bodyReader := bytes.NewBuffer(j) //transforma jsonul in io.Reader

	res, err := http.NewRequest(http.MethodPost, callbackurl, bodyReader)
	if err != nil {
		log.Printf("Error making request back: %v", err)
		return err
	}
	res.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Error executing request back: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Print("Handshake ok")
	} else {
		log.Print("Handshake failed")
		return fmt.Errorf("Error handshake")
	}

	return nil

}
