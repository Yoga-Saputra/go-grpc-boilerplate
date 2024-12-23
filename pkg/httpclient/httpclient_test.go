package httpclient

import (
	"fmt"
	"strings"
	"testing"
)

func createNew() *Config {
	return New(Config{BaseUrl: "https://httpbin.org"})
}

func get(c *Config, e string) error {
	resp, err := c.Get(e)
	fmt.Println("Result:", resp.GetJSONString())

	return err
}

func post(c *Config, e string, p ...map[string]interface{}) error {
	resp, err := c.Post(e, p...)
	fmt.Println("Result:", resp.GetJSONString())

	return err
}

func Test_Get(t *testing.T) {
	if err := get(createNew(), "/get"); err != nil {
		t.Errorf("httpclient.Get() got error '%v', expected '%v'", err.Error(), nil)
	}
}

func Test_Timeout(t *testing.T) {
	c := createNew()
	c.SetTimeout(2)
	if err := get(c, "/delay/4"); err != nil {
		if !strings.Contains(err.Error(), "context deadline") {
			t.Errorf("httpclient.Get() with time out got error '%v', expected '%v'", err.Error(), nil)
		}
	}
}

func Test_BearerToken(t *testing.T) {
	c := createNew()
	c.Token("asdfasdf")
	if err := get(c, "/bearer"); err != nil {
		t.Errorf("httpclient.Get() with time bearer token got error '%v', expected '%v'", err.Error(), nil)
	}
}

func Test_Post(t *testing.T) {
	if err := post(createNew(), "/anything"); err != nil {
		t.Errorf("httpclient.Post() got error '%v', expected '%v'", err.Error(), nil)
	}
}
