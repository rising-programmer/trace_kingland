package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query callback representing the query of a chaincode
func (t *PeopleChainCode) richQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("输入参数数量应=1！")
	}

	fmt.Println("querystring:", args[0])
	rqi, err := stub.GetQueryResult(args[0])
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
		//table, keys, err := stub.SplitCompositeKey(queryResponse.Key)
		_, _, err = stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if i > 1 {
			buffer.WriteString(",")
		}
		/*buffer.WriteString("{\"Table\":")
		buffer.WriteString("\"")
		buffer.WriteString(table)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(keys[0])
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is*/
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
	}
	buffer.WriteString("]")

	if i == 0 {
		return shim.Success(nil)
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
