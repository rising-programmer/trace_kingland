package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AuditHistory struct {
	TxId     string
	Value    []byte
	TxTime   int64
	IsDelete bool
}

// query for history
func (t *PeopleChainCode) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var historyList []AuditHistory
	var obj []byte

	if len(args) != 2 {
		return shim.Error("输入参数数量应=2！")
	}

	var tableName string = args[0] //表名
	var key string = args[1]       //主键

	//复合主键
	compositeKey, _ := stub.CreateCompositeKey(tableName, []string{key})

	fmt.Println("query for history : type: %s, Key : %s .", tableName, compositeKey)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(compositeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryresult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		txID := queryresult.GetTxId()
		historicValue := queryresult.GetValue()
		txUnixTime := queryresult.GetTimestamp().GetSeconds()
		isDelete := queryresult.GetIsDelete()
		fmt.Println(txUnixTime)
		//docker时区问题，时间输出为unix时间戳

		if isDelete {
			obj = nil
		} else {
			obj = historicValue
		}

		var auditHistory AuditHistory
		auditHistory.TxId = txID
		auditHistory.Value = obj
		auditHistory.TxTime = txUnixTime
		auditHistory.IsDelete = isDelete

		historyList = append(historyList, auditHistory)

	}

	var listAsBytes []byte

	listAsBytes, err = json.Marshal(historyList)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(listAsBytes)
}
