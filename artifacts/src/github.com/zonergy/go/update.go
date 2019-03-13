package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

func (t *PeopleChainCode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var id string = args[0]        //数据表简单主键
	var argStr string = args[1]    //json字符串
	bytes := []byte(argStr)        //存入链中的内容
	var err error

	//将json数据写入区块链
	err = stub.PutState(id, bytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
