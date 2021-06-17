package omni

import (
	"bytes"
	"fmt"
	"net/http"
)

type ConnConfig struct {
	// The format is <host|ip:[port]>
	Host string
	// Give a username for Authentication
	User string
	// Give a password for Authentication
	Pass string
}

type Client struct {
	connConfig *ConnConfig
}

func NewClient(connConfig *ConnConfig) *Client {
	return &Client{connConfig}
}

func (c *Client) Call(result interface{}, method string, args ...interface{}) error {
	url := fmt.Sprintf("http://%s/", c.connConfig.Host)
	params := make([]interface{}, 0)
	params = append(params, args...)

	message, err := EncodeClientRequest(method, params)
	if err != nil {
		//log.Printf("%s", err,",EncodeClientRequest")
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		//log.Printf("%s", err,",NewRequest")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.connConfig.User, c.connConfig.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Printf("Error in sending request to %s. %s", url, err)
		return err
	}
	defer resp.Body.Close()

	err = DecodeClientResponse(resp.Body, &result)
	if err != nil {
		//log.Printf("Couldn't decode response. %s", err)
		return err
	}
	return nil
}

func (c *Client) CallWithoutResponse(method string, args ...interface{}) error {
	url := fmt.Sprintf("http://%s/", c.connConfig.Host)
	params := make([]interface{}, 0)
	params = append(params, args...)

	message, err := EncodeClientRequest(method, params)
	if err != nil {
		//log.Printf("%s", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		//log.Printf("%s", err,"NewRequest 2")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.connConfig.User, c.connConfig.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Printf("Error in sending request to %s. %s", url, err)
		return err
	}
	defer resp.Body.Close()
	
	return nil
}
