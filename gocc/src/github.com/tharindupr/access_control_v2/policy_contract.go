package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)




//creating a subject asset
func (s *SmartContract) createPolicy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	// if len(args) != 5 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 5")
	// }
	

	id, _ := cid.GetID(APIstub)
	
	var sattributes = map[string]string{};
	json.Unmarshal([]byte(args[1]), &sattributes)

	var oattributes = map[string]string{};
	json.Unmarshal([]byte(args[2]), &oattributes)

	var rule = []Rule{};
	json.Unmarshal([]byte(args[3]), &rule)

	var policy = Policy{UserName: args[0], CID: id, SubjectAttributes: sattributes, ObjectAttributes: oattributes, Rules: rule}

	policyAsBytes, _ := json.Marshal(policy)
	APIstub.PutState(id, policyAsBytes)

	logger.Infof("Successfully Added")
	return shim.Success(policyAsBytes)
}


// query Policy
func (s *SmartContract) queryPolicy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	subjectAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(subjectAsBytes)
}


func (t *SmartContract) getHistoryForPolicy(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetID := args[0]
	logger.Infof("searching for object id %s", args[0])
	resultsIterator, err := stub.GetHistoryForKey(assetID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		logger.Infof(string(response.TxId))
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	logger.Infof("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


// adding a rule to a rule 
func (s *SmartContract) addRule(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	policyAsBytes, _ := APIstub.GetState(args[0])
	policy:= Policy{}

	var rule = Rule{};
	json.Unmarshal([]byte(args[1]), &rule)

	json.Unmarshal(policyAsBytes, &policy)
	policy.Rules = append(policy.Rules, rule)
	

	policyAsBytes, _ = json.Marshal(policy)
	APIstub.PutState(args[0], policyAsBytes)

	return shim.Success(policyAsBytes)
}