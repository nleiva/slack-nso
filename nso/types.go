package nso

// Router is a IOS XR tail-f router (tailf-ned-cisco-ios-xr:router"). Only static routing config for now
type Router struct {
	Static static `json:"tailf-ned-cisco-ios-xr:static"`
}

type route struct {
	Net     string `json:"net"`
	Address string `json:"address"`
}

type uni struct {
	Unicast struct {
		Routes []route `json:"routes-ip"`
	} `json:"unicast"`
}

type static struct {
	AddressFamily struct {
		IPv4 uni `json:"ipv4"`
		IPv6 uni `json:"ipv6"`
	} `json:"address-family"`
}
