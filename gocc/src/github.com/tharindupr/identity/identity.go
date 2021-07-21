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
	AssetID   string `json:"assetid"`
	Type  string `json:"type"`
	Attributes map[string]string `json:"attributes"`
	CID  string `json:"cid"`
	Status string `json:"status"`
}


// Private Assets Structure
type PrivateAsset struct {
	AssetID   string `json:"assetid"`
	Type  string `json:"type"`
	CID  string `json:"cid"`
}


// Private Assets Structure
type assetPrivateDetails struct {
	AssetID   string `json:"assetid"`
	Attributes map[string]string `json:"attributes"`
	//FixedAttribute1 string `json:"fixedattribute1"`
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
		case "createAsset":
			logger.Infof("Case passed")
			return s.createAsset(APIstub, args)
		case "getAsset":
			return s.getAsset(APIstub, args)
		case "traceAsset":
			return s.traceAsset(APIstub, args)
		case "addAttribute":
			return s.addAttribute(APIstub, args)
		case "createPrivateAsset":
			return s.createPrivateAsset(APIstub, args)
		case "getPrivateAsset":
			return s.getPrivateAsset(APIstub, args)
		case "queryPrivateDataHash":
			return s.queryPrivateDataHash(APIstub, args)
		case "updateAssetStatus":
			return s.updateAssetStatus(APIstub, args)
			
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
func (s *SmartContract) createAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	logger.Infof("In fucction createAsset")
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	//checking whether the key exists
	assetAsBytes, _ := APIstub.GetState(args[0])
	if assetAsBytes != nil {
		return shim.Error("Key Exist Already")
	}


	logger.Infof(args[3])
	//creating an object from the attribute array
	var attributes = map[string]string{};
	json.Unmarshal([]byte(args[3]), &attributes)

	var asset = Asset{AssetID: args[0], Type: args[1], Status: args[2], Attributes: attributes, CID: "NULL"}

	logger.Infof("Saving")
	//logger.Infof(subject.Attributes)
	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(assetAsBytes)
}

//updating a digital energy certificate
func (s *SmartContract) updateAssetStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}


	// ABAC
	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err !=nil {
		return shim.Error("Error retriving user attributes")
	}
	
	if !ok {
		//The client identity does not possess the attributes
		return shim.Error("The client identity does not possess the attributes")
	}

	if val != "admin" {
		fmt.Println("Attribute role : " + val)
		return shim.Error("Insufficient permisions to update the status of the Asset. Only Admins can change status")
	}

	//checking whether the key exists
	assetAsBytes, _ := APIstub.GetState(args[0])
	if assetAsBytes == nil {
		return shim.Error("Key Doesn't Exist")
	}


	asset := Asset{}

	json.Unmarshal(assetAsBytes, &asset)
	asset.Status = args[1]

	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(assetAsBytes)
}


//query subject
func (s *SmartContract) getAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(assetAsBytes)
}


func (t *SmartContract) traceAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

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

func (s *SmartContract) createPrivateAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type assetTransientInput struct {
		AssetID  string `json:"assetid"` //the fieldtags are needed to keep case from bouncing around
		Type string `json:"type"`
		Attributes  map[string]string `json:"attributes"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	assetDataAsBytes, ok := transMap["asset"]
	if !ok {
		return shim.Error("asset must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(assetDataAsBytes))

	if len(assetDataAsBytes) == 0 {
		return shim.Error("333333 -asset value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var assetInput assetTransientInput
	err = json.Unmarshal(assetDataAsBytes, &assetInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(assetDataAsBytes) + "Error is : " + err.Error())
	}

	logger.Infof("3333")

	// if len(carInput.Key) == 0 {
	// 	return shim.Error("name field must be a non-empty string")
	// }
	// if len(carInput.Make) == 0 {
	// 	return shim.Error("color field must be a non-empty string")
	// }
	// if len(carInput.Model) == 0 {
	// 	return shim.Error("model field must be a non-empty string")
	// }
	// if len(carInput.Color) == 0 {
	// 	return shim.Error("color field must be a non-empty string")
	// }
	// if len(carInput.Owner) == 0 {
	// 	return shim.Error("owner field must be a non-empty string")
	// }
	// if len(carInput.Price) == 0 {
	// 	return shim.Error("price field must be a non-empty string")
	// }

	logger.Infof("444444")

	// ==== Check if the asset already exists ====
	assetAsBytes, err := APIstub.GetPrivateData("collectionAssets", assetInput.AssetID)
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes != nil {
		fmt.Println("This asset already exists: " + assetInput.AssetID)
		return shim.Error("This asset already exists: " + assetInput.AssetID)
	}

	logger.Infof("55555")

	var asset = PrivateAsset{AssetID: assetInput.AssetID, Type: assetInput.Type, CID: "NULL"}

	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutPrivateData("collectionAssets", assetInput.AssetID, assetAsBytes)
	if err != nil {
		logger.Infof("6666666")
		return shim.Error(err.Error())
	}

	assetPrivateDetails := &assetPrivateDetails{AssetID: assetInput.AssetID, Attributes: assetInput.Attributes}

	assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionAssetPrivateDetails", assetInput.AssetID, assetPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(assetAsBytes)
}


func (s *SmartContract) getPrivateAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// collectionCars, collectionCarPrivateDetails, _implicit_org_Org1MSP, _implicit_org_Org2MSP
	assetAsBytes, err := APIstub.GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if assetAsBytes == nil {
		jsonResp := "{\"Error\":\"Asset private details does not exist: " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(assetAsBytes)
}

func (s *SmartContract) queryPrivateDataHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	assetAsBytes, _ := APIstub.GetPrivateDataHash(args[0], args[1])

	return shim.Success(assetAsBytes)
}