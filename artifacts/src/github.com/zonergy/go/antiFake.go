/**
 * 防伪验证
 */
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

func (t *PeopleChainCode) antiFake(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//输入参数校验
	if len(args) != 3 {
		return shim.Error("输入参数数量应=3！")
	}

	var tableName string = args[0] //表名
	var id string = args[1]        //数据表简单主键
	var argStr string = args[2]    //json字符串
	bytes := []byte(argStr)        //存入链中的内容
	var requestAntiFakeCode string
	var err error
	var paramMap map[string]string

	var m map[string]interface{}
	err = json.Unmarshal(bytes,&paramMap)
	if err != nil{
		return shim.Error("{\"Error\":\"jsonStr unMarshal failed . source = "+string(bytes)+"\"}")
	}
	requestAntiFakeCode,ok := paramMap["Data"]
	if !ok  {
		return shim.Error("{\"Error\":\"requested param securityCode is not exist\"}")
	}

	//拼接为复合主键
	//key = tableName + "-" + id
	compositeKey, _ := stub.CreateCompositeKey(tableName, []string{id})
	bytes, err = stub.GetState(compositeKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + compositeKey + "\"}"
		return shim.Error(jsonResp)
	}
	if bytes == nil{
	 	return shim.Error("{\"Error\":\"entity is not exist \"}")
	}
	//JSON To Map
	err = json.Unmarshal(bytes,&m)
	if err != nil{
		return shim.Error("{\"Error\":\"jsonToMap convert failed \"}")
	}
	//verify anti-code
	data,ok := m["Data"].(map[string]interface{})
	if !ok {
		return shim.Error("{\"Error\":\"key `Data` is not exist\"}")
	}
	securityCode,ok := data["securityCode"]
	
	if !ok {
		return shim.Error("{\"Error\":\"securityCode is not exist\"}")
	}
	if securityCode != requestAntiFakeCode{
		return shim.Error("{\"Error\":\"securityCode is not matched\"}")
	}

	compositeKey_securityCode, _ := stub.CreateCompositeKey(tableName, []string{id,securityCode.(string)})
	bytes, err = stub.GetState(compositeKey_securityCode)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + compositeKey_securityCode + "\"}")
	}
	if bytes != nil{
	 	return shim.Error("{\"Error\":\"securityCode is used\"}")
	}else{
		//onchain securityCode
		 err = stub.PutState(compositeKey_securityCode, []byte(securityCode.(string)))
		 if err != nil {
		 	return shim.Error(err.Error())
		 }
	}

	//If TxId exist,then return it.
	txId,ok := m["TxId"].(string)
	if ok{
		return shim.Success([]byte(txId))
	}else{
		return shim.Error("{\"Error\":\"TxId can not be empty \"}")
	}
}
