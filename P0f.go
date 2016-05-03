// Package P0f is trying to create an API to pass request to unix socket of
// P0f3 http://lcamtuf.coredump.cx/p0f3/
//
// The goal is to be able to query an ip to P0f to retrieve informations about
// the wanted host
//
// Documentation of the API is located here : http://lcamtuf.coredump.cx/p0f3/README
//
package P0f

//go:generate stringer -type=osMatchQType,badSwType

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// MAGIC_QUERY is the magic dword that initiate a request to the unix socket
const magicQuery int32 = 0x50304601

// MAGIC_RESPONSE is the magic dword that start an unix socket response from p0f
const magicResponse int32 = 0x50304602

// ResponseStatusType is the type of response in regard to the query done.
type responseStatusType int32

const (
	BadQuery responseStatusType = 0
	OK                          = 0x10
	NoMatch                     = 0x20
)

// OsMatchType is an estimation of the quality of the detection of the Host.
type osMatchQType int8

const (
	Normal osMatchQType = iota
	Fuzzy
	Generic
	Both
)

// BadSwType is representing the confidence we can have in the user-Agent value
type badSwType int8

const (
	Nul badSwType = iota
	OsMismatch
	Mismatch
)

type response struct {
	MagicWord      int32
	Status         responseStatusType
	FirstSeen      int32
	LastSeen       int32
	TotalConn      int32
	UptimeMin      int32
	UpModDays      int32
	LastNAT        int32
	LastChg        int32
	Distance       int16
	BadSw          badSwType
	OsMatchQuality osMatchQType
	OsName         [32]byte
	OsFlavor       [32]byte
	HttpName       [32]byte
	HttpFlavor     [32]byte
	LinkType       [32]byte
	Language       [32]byte
}

// Response struct is describing all informations retrieved from P0f unix socket.
type Response struct {
	FirstSeen      time.Time
	LastSeen       time.Time
	TotalConn      int
	UptimeMin      int
	UpModDays      int
	LastNAT        time.Time
	LastChg        time.Time
	Distance       int
	BadSw          badSwType
	OsMatchQuality osMatchQType
	OsName         string
	OsFlavor       string
	HttpName       string
	HttpFlavor     string
	LinkType       string
	Language       string
}

func (r *Response) String() string {
	return fmt.Sprintln(
		"FirstSeen:", r.FirstSeen,
		"\nLastSeen:", r.LastSeen,
		"\nBadSw:", r.BadSw,
		"\nOsMatchQuality:", r.OsMatchQuality,
		"\nOsName:", r.OsName,
		"\nOsFlavor:", r.OsFlavor,
		"\nHttpName:", r.HttpName,
		"\nLinkType:", r.LinkType,
		"\nLanguage:", r.Language,
	)
}

// AddressType represent code value of ipv4/ipv6
type addressType int8

const (
	IPv4 addressType = 4
	IPv6             = 6
)

type request struct {
	MagicWord   int32
	AddressType addressType
	Address     [16]byte
}

func newRequest(ip net.IP) (*request, error) {
	var req = new(request)
	req.MagicWord = magicQuery

	tmp := ip.To4()
	if tmp == nil {
		req.AddressType = IPv6
		copy(req.Address[:], []byte(ip))
	} else {
		req.AddressType = IPv4
		for index, b := range []byte(tmp) {
			req.Address[index] = b
		}
	}
	return req, nil
}

// P0f is the type that handle connection to the p0f unix socket
type P0f struct {
	conn net.Conn
}

// New return a P0f Object. unixSocket must be the socket defined by -s option of p0f.
// P0f object would be used to query your local p0f instance to obtain data
// about an IPAddress
func New(unixSocket string) (*P0f, error) {
	c, err := net.Dial("unix", unixSocket)
	if err != nil {
		return nil, err
	}

	return &P0f{conn: c}, nil
}

// GetAddrInfo uses IP address in both IPv4 and IPv6 formats
// It returns an P0f.Response type that describe the queried IP.
// In case of error, a nil type is returned
func (p *P0f) GetAddrInfo(sip string) (*Response, error) {
	ip := net.ParseIP(sip)
	if ip == nil {
		return nil, fmt.Errorf("Couldn't parse %v as an IP address", sip)
	}
	return p.GetIPInfo(ip)
}

// GetIPInfo uses net.IP format and return a P0f.Response
// corresponding to informations about the queried IP address
func (p *P0f) GetIPInfo(ip net.IP) (*Response, error) {
	var req *request
	var err error
	req, err = newRequest(ip)
	if err != nil {
		return nil, err
	}

	err = binary.Write(p.conn, binary.LittleEndian, req)
	if err != nil {
		return nil, err
	}

	var resp response
	err = binary.Read(p.conn, binary.LittleEndian, &resp)
	//fmt.Println("Status of response: ", resp.Status)
	//fmt.Println("Response:", resp)

	switch resp.Status {
	case BadQuery:
		return nil, fmt.Errorf("Bad query")
	case OK:
		return &Response{
			FirstSeen:      time.Unix(int64(resp.FirstSeen), 0),
			LastSeen:       time.Unix(int64(resp.LastSeen), 0),
			TotalConn:      int(resp.TotalConn),
			UptimeMin:      int(resp.UptimeMin),
			UpModDays:      int(resp.UpModDays),
			LastNAT:        time.Unix(int64(resp.LastNAT), 0),
			LastChg:        time.Unix(int64(resp.LastChg), 0),
			Distance:       int(resp.Distance),
			BadSw:          resp.BadSw,
			OsMatchQuality: resp.OsMatchQuality,
			OsName:         string(resp.OsName[:]),
			OsFlavor:       string(resp.OsFlavor[:]),
			HttpName:       string(resp.HttpName[:]),
			HttpFlavor:     string(resp.HttpFlavor[:]),
			LinkType:       string(resp.LinkType[:]),
			Language:       string(resp.Language[:]),
		}, nil
	case NoMatch:
		return nil, fmt.Errorf("No match")
	}

	return nil, fmt.Errorf("Unrecognized response")
}
