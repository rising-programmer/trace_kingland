package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query callback representing the query of a chaincode
func (t *PeopleChainCode) export(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("输入参数数量应=1！")
	}

	var table string = args[0] //表名

	rqi, err := stub.GetStateByPartialCompositeKey(table, []string{})
	if err != nil {
		shim.Error(err.Error())
	}
	defer rqi.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	i := 0
	fmt.Println("Running loop")
	for rqi.HasNext() {
		i++
		queryResponse, err := rqi.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		_, _, err = stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		if i > 1 {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
	}
	buffer.WriteString("]")

	if i == 0 {
		return shim.Success(nil)
	}

	return shim.Success(buffer.Bytes())
}
