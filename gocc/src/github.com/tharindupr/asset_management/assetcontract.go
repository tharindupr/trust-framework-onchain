package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Assets Structure
type Asset struct {
	ID   string `json:"id"`
	Type  string `json:"type"`
	Attributes map[string]string `json:"attributes"`
	CID  string `json:"cid"`
}


var logger = flogging.MustGetLogger("subject_cc")


// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "createSubject":
		return s.createSubject(APIstub, args)
	case "createObject":
		return s.createObject(APIstub, args)
	case "queryObject":
		return s.queryObject(APIstub, args)
	case "querySubject":
		return s.querySubject(APIstub, args)
	case "queryAssetHistory":
		return s.getHistoryForAsset(APIstub, args)
	case "addAttribute":
		return s.addAttribute(APIstub, args)
	case "createCar":
		return s.createObject(APIstub, args)
	}
	
	return shim.Error("Invoke Function Not Success.")

}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}


//creating a subject asset
func (s *SmartContract) createSubject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//creating an object from the attribute array
	var attributes = map[string]string{};
	json.Unmarshal([]byte(args[2]), &attributes)

	var subject = Asset{ID: args[0], Type: args[1], Attributes: attributes, CID: "NULL"}

	//logger.Infof(subject.Attributes)
	subjectAsBytes, _ := json.Marshal(subject)
	APIstub.PutState(args[0], subjectAsBytes)

	return shim.Success(subjectAsBytes)
}


//creating a object asset
func (s *SmartContract) createObject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//getting ID for the client which create the Object
	id, _ := cid.GetID(APIstub)
	args = append(args, id)


	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	//creating an object from the attribute array
	var attributes =  map[string]string{};
	json.Unmarshal([]byte(args[2]), &attributes)

	var object = Asset{ID: args[0], Type: args[1], Attributes: attributes, CID: id}

	objectAsBytes, _ := json.Marshal(object)
	APIstub.PutState(args[0], objectAsBytes)

	return shim.Success(objectAsBytes)
}

//query subject
func (s *SmartContract) querySubject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	subjectAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(subjectAsBytes)
}


func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

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

//query object
func (s *SmartContract) queryObject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	subjectAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(subjectAsBytes)
}

//add attributes
func (s *SmartContract) addAttribute(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	assetAsBytes, _ := APIstub.GetState(args[0])
	asset:= Asset{}

	json.Unmarshal(assetAsBytes, &asset)
	asset.Attributes[args[1]]= args[2]

	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(assetAsBytes)
}

