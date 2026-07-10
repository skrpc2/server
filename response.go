package rysrv

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fxamacker/cbor/v2"
	"github.com/golang/protobuf/proto"
)

// SetError writes JSON-RPC response with error.
//
// It overwrites previous calls of SetResult and SetError.
func SetError(c *app.RequestContext, err error) {
	args := &Base{}
	args.Err = err.Error()

	b, err := proto.Marshal(args)
	if err != nil {
		fmt.Println("SetError args.MarshalBinary = ", err.Error())
		return
	}
	c.Response.SetBody(b)
}

func SetResult(c *app.RequestContext, result interface{}) {
	args := &Base{}
	args.Err = ""

	if vv, ok := result.(string); ok {
		args.Data = []byte(vv)
	} else if vv, ok := result.([]byte); ok {
		args.Data = vv
	} else {
		b1, err1 := cbor.Marshal(result)
		if err1 != nil {
			fmt.Println("cbor.Marshal err = ", err1.Error())
			return
		}
		args.Data = b1
	}

	b2, err2 := proto.Marshal(args)
	if err2 != nil {
		fmt.Println("SetResult args.MarshalBinary = ", err2.Error())
		return
	}
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.SetBody(b2)
}
