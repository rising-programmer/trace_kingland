package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query callback representing the query of a chaincode
func (t *PeopleChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	var key string = args[0]   //主键
	var objByte []byte         //存入链中的内容

	objByte, err = stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	if objByte == nil {
		jsonResp := "{\"Error\":\"Nil context for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(objByte)
}
