package P0f

// This stringer is made manually because types doesn't use iota and uses custom values
// this is not currently supported in go generate code

import "fmt"

const _AddressType_name = "IPv4IPv6"

var _AddressType_index = [...]uint8{0, 4, 8}

func (ii addressType) String() string {
	var i int
	switch ii {
	case 4:
		i = 0
	case 6:
		i = 1
	default:
		i = int(ii)
	}
	if i < 0 || i >= len(_AddressType_index)-1 {
		return fmt.Sprintf("AddressType(%d)", ii)
	}
	return _AddressType_name[_AddressType_index[i]:_AddressType_index[i+1]]
}

const _ResponseStatusType_name = "BadQueryOKNoMatch"

var _ResponseStatusType_index = [...]uint8{0, 8, 10, 17}

func (ii responseStatusType) String() string {
	var i int
	switch ii {
	case 0x00:
		i = 0
	case 0x10:
		i = 1
	case 0x20:
		i = 2
	default:
		i = int(ii)
	}
	if i < 0 || i >= len(_ResponseStatusType_index)-1 {
		return fmt.Sprintf("ResponseStatusType(%d)", ii)
	}
	return _ResponseStatusType_name[_ResponseStatusType_index[i]:_ResponseStatusType_index[i+1]]
}
