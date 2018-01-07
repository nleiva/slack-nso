package nso

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
	return
}

func decodeJSON(v interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(v)
}

func generateStatic(n string, a string) (string, error) {
	r := new(Router)
	rt := route{Net: n, Address: a}

	ip, _, err := net.ParseCIDR(n)
	if err != nil {
		return "", err
	}
	// To4 converts the IPv4 address ip to a 4-byte representation. If ip is not an IPv4 address, To4 returns nil.
	if ip.To4() != nil {
		r.Static.AddressFamily.IPv4.Unicast.Routes = append(r.Static.AddressFamily.IPv4.Unicast.Routes, rt)
	} else {
		r.Static.AddressFamily.IPv6.Unicast.Routes = append(r.Static.AddressFamily.IPv6.Unicast.Routes, rt)
	}

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func readStatic(r Router) {
	// Just examples for now
	fmt.Printf("data4: %v\n", r.Static.AddressFamily.IPv4.Unicast.Routes[0])
	fmt.Printf("data4: %v\n", r.Static.AddressFamily.IPv4.Unicast.Routes[1].Net)
	fmt.Printf("data4: %v\n", r.Static.AddressFamily.IPv4.Unicast.Routes[1].Address)
	fmt.Printf("data6: %v\n", r.Static.AddressFamily.IPv6.Unicast.Routes[0])
}
