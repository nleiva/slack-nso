package nso

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Server is a NSO server
type Server struct {
	Addr *url.URL
}

// Configurator configures Network Elements
type Configurator interface {
	StaticRoute(msg []string)
}

// StaticRoute configures a static route via NSO
func (s *Server) StaticRoute(msg []string, device string) {
	var netClient = &http.Client{
		Timeout: time.Second * 20,
	}

	//req, err := fullConfig(u, device")
	//req, err := interfaceConfig(u, device)
	//req, err := routerConfig(u, device)
	//req, err := syncFrom(u, device)
	//config, err := generateStatic("191.0.0.0/8", "10.87.89.1")
	//config, err := generateStatic("2001:425::/32", "2001:420:2cff:1204::1")
	config, err := generateStatic(msg[1], msg[2])
	checkErr(err)

	req, err := setRouterConfig(s.Addr, device, "static", config)
	checkErr(err)

	resp, err := netClient.Do(req)
	checkErr(err)
	defer resp.Body.Close()

	// Read JSON data
	/* 	data := new(Router)
	   	err = decodeJSON(data, resp.Body)
	   	checkErr(err)
	   	fmt.Printf("%v\n", data)
	*/

	contents, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	fmt.Printf("%s\n", string(contents))
}

func fullConfig(u *url.URL, d string) (req *http.Request, err error) {
	// "http://admin:admin@mrstn-nso.cisco.com:8080/api/config/devices/device/mrstn-5501-1.cisco.com/config?deep=true"
	u.Path = "api/config/devices/device/" + d + "/config"
	// All the details
	q := u.Query()
	q.Set("deep", "true")
	u.RawQuery = q.Encode()
	// Request
	req, err = http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", "application/vnd.yang.data+json")
	return req, err
}

func interfaceConfig(u *url.URL, d string) (req *http.Request, err error) {
	u.Path = "api/running/devices/device/" + d + "/config/interface/TenGigE"
	// All the details
	q := u.Query()
	q.Set("deep", "true")
	u.RawQuery = q.Encode()
	// Request
	req, err = http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", "application/vnd.yang.collection+json")
	return req, err
}

func routerConfig(u *url.URL, d string) (req *http.Request, err error) {
	// static, isis, bgp, etc...
	p := "static"
	u.Path = "api/running/devices/device/" + d + "/config/router/" + p
	// All the details
	q := u.Query()
	q.Set("deep", "true")
	u.RawQuery = q.Encode()
	// Request
	req, err = http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Accept", "application/vnd.yang.data+json")
	return req, err
}

func setRouterConfig(u *url.URL, d string, p string, c string) (req *http.Request, err error) {
	// PUT, POST -> REPLACE. PATH -> MERGE
	// p = static, isis, bgp, etc...
	u.Path = "api/running/devices/device/" + d + "/config/router/" + p
	// Request
	req, err = http.NewRequest("PATCH", u.String(), strings.NewReader(c))
	req.Header.Add("Content-Type", "application/vnd.yang.data+json")
	req.Header.Add("Accept", "application/vnd.yang.data+json")
	return req, err
}

func syncFrom(u *url.URL, d string) (req *http.Request, err error) {
	u.Path = "api/running/devices/device/" + d + "/_operations/sync-from"
	// Request
	req, err = http.NewRequest("POST", u.String(), nil)
	req.Header.Add("Content-Type", "application/vnd.yang.operation+json")
	req.Header.Add("Accept", "application/vnd.yang.operation+json")
	return req, err
}
