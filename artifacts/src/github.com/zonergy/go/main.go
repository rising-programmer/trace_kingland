package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PeopleChainCode struct {
}

func (t *PeopleChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}
func (t *PeopleChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("function: %s,args: %s\n", function, args)
	if function == "update" {
		return t.update(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "history" {
		return t.history(stub, args)
	} else if function == "export" {
		return t.export(stub, args)
	} else if function == "richQuery" {
		return t.richQuery(stub, args)
	} else if function == "antiFake" {
		return t.antiFake(stub,args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"export\" \"query\" \"history\" \"delete\"")
}

func main() {
	err := shim.Start(new(PeopleChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
