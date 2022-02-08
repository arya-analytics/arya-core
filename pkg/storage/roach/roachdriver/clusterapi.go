package roachdriver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ClusterAPI struct {
	Username     string
	Password     string
	Host         string
	Port         int
	sessionToken string
}

func (c *ClusterAPI) Connect() {
	client := &http.Client{}
	req, err := http.NewRequest("POST",
		c.apiURL()+"login/?username=root&password=testpass", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	c.parseSessionToken(resp.Body)
}

func (c *ClusterAPI) Nodes() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.apiURL()+"nodes/", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set(sessionTokenKey, c.sessionToken)
	log.Info(req.Header.Get(sessionTokenKey))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Info(resp.Request.Header)
	body, err := ioutil.ReadAll(resp.Body)
	log.Info(string(body))
	time.Sleep(10000 * time.Second)
}

func (c *ClusterAPI) apiURL() string {
	return fmt.Sprintf("https://%s:%v/api/v2/", c.Host, c.Port)
}

func (c *ClusterAPI) parseSessionToken(respBody io.Reader) {
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Fatalln(err)
	}
	bodyStr := string(body)
	token := strings.Split(strings.Split(bodyStr, ":")[1], "\"}")[0][1:]
	c.sessionToken = token
}

const (
	sessionTokenKey = "X-Cockroach-API-Session"
)

func (c *ClusterAPI) addTokenHeader(req *http.Request) {
	req.Header.Add(sessionTokenKey, c.sessionToken)

}
