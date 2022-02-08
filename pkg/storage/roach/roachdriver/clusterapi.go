package roachdriver

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const protocol = "https://"

type clusterAPIEndpoint string

const (
	clusterAPIEndpointBase  clusterAPIEndpoint = "/api/v2/"
	clusterAPIEndpointLogin clusterAPIEndpoint = "login/"
	clusterAPIEndpointNodes clusterAPIEndpoint = "nodes/"
)

type ClusterAPI struct {
	Username     string
	Password     string
	Host         string
	Port         int
	sessionToken string
	_client      *http.Client
}

func (c *ClusterAPI) authParamString() string {
	return fmt.Sprintf("?username=%s&password=%s", c.Username, c.Password)
}

func (c *ClusterAPI) client() *http.Client {
	if c._client == nil {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c._client = &http.Client{}
	}
	return c._client
}

func (c *ClusterAPI) Connect() error {
	req := c.buildPOSTRequest(clusterAPIEndpointLogin, c.authParamString())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}
	return c.parseSessionToken(resp.Body)
}

func (c *ClusterAPI) Nodes() (ClusterAPINodes, error) {
	resp, err := c.doGETRequest(clusterAPIEndpointNodes, "")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jsonBody := map[string]ClusterAPINodes{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return nil, err
	}
	return jsonBody["nodes"], err
}

func (c *ClusterAPI) doGETRequest(ep clusterAPIEndpoint, ext string) (*http.Response,
	error) {
	return c.client().Do(c.buildGETRequest(ep, ext))
}

func (c *ClusterAPI) buildGETRequest(ep clusterAPIEndpoint, ext string) *http.Request {
	return c.buildRequest("GET", ep, ext)
}

func (c *ClusterAPI) buildPOSTRequest(ep clusterAPIEndpoint,
	ext string) *http.Request {
	return c.buildRequest("POST", ep, ext)

}

func (c *ClusterAPI) buildRequest(method string, ep clusterAPIEndpoint,
	ext string) *http.Request {
	u := c.buildURL(ep, ext)
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		panic(err)
	}
	if c.sessionToken != "" {
		c.addSessionTokenHeader(req)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (c *ClusterAPI) buildURL(ep clusterAPIEndpoint, ext string) *url.URL {
	u, err := url.Parse(protocol + c.Host + ":" + strconv.Itoa(c.
		Port) + string(clusterAPIEndpointBase) + string(ep) + ext)
	// because we aren't accepting any outside input for urls,
	// a failure to parse is a programmatic error
	if err != nil {
		panic(err)
	}
	return u
}

func (c *ClusterAPI) parseSessionToken(respBody io.Reader) error {
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	token := strings.Split(strings.Split(bodyStr, ":")[1], "\"}")[0][1:]
	c.sessionToken = token
	return nil
}

const sessionTokenKey = "X-Cockroach-API-Session"

func (c *ClusterAPI) addSessionTokenHeader(req *http.Request) {
	req.Header.Add(sessionTokenKey, c.sessionToken)

}

type ClusterAPINodes []ClusterAPINode

type ClusterAPINodeMetrics map[string]int

type ClusterAPINode struct {
	ID      int `json:"node_id"`
	Metrics ClusterAPINodeMetrics
}
