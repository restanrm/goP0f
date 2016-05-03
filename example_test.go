package P0f_test

import (
	"fmt"

	"github.com/restanrm/goP0f"
)

func ExampleP0f_GetAddrInfo() {
	pof, err := P0f.New("/dev/shm/test") // unix socket file of p0f, created with command `p0f -s /dev/shm/test`
	if err != nil {
		fmt.Println(err)
		return
	}

	var resp *P0f.Response
	resp, err = pof.GetAddrInfo("192.168.1.1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
