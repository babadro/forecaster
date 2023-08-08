package main

import (
	"fmt"

	"github.com/babadro/forecaster/cmd/temp/proto/callbackdata"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/golang/protobuf/proto"
)

func main() {
	data := &callbackdata.CallbackData{
		Field1:  helpers.Ptr[int32](127),
		Field2:  helpers.Ptr[int32](127),
		Field3:  helpers.Ptr[int32](127),
		Field4:  helpers.Ptr[int32](127),
		Field5:  helpers.Ptr[int32](127),
		Field6:  helpers.Ptr[int32](127),
		Field7:  helpers.Ptr[int32](127),
		Field8:  helpers.Ptr[int32](127),
		Field9:  helpers.Ptr[int32](127),
		Field10: helpers.Ptr[int32](127),
		Field11: helpers.Ptr[int32](127),
		Field12: helpers.Ptr[int32](127),
		Field13: helpers.Ptr[int32](127),
		Field14: helpers.Ptr[int32](127),
		Field15: helpers.Ptr[int32](127),
	}

	binaryData, err := proto.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(binaryData)

	fmt.Println("len(binaryData):", len(binaryData))
}
