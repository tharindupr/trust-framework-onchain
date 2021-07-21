package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	//"reflect"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	//"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}


// Digital Energy Certificate Structure
type DEC struct {
	DECID   string `json:"decid"`
	BuildingID  string `json:"buildingid"`
	//Attributes map[string]string `json:"attributes"`
	Status string `json:"status"`
	CID  string `json:"cid"`
	BuildingCategory string `buildingcategory`
	FloorArea string `floorarea`
	HoursOfOccupancy float64 `hoursofoccupancy`
	EnergyConsumption float64 `energyconsumption`
	MeterStartDate string `meterstartdate`
	MetereEndDate string `meterenddate`
	Grade string `grading`

}

var logger = flogging.MustGetLogger("subject_cc")


// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("Chiancode : accesscontrolcontract initiated")
	return shim.Success(nil)
}


// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
		case "createDEC":
			return s.createDEC(APIstub, args)
		case "updateDEC":
			return s.updateDEC(APIstub, args)
		case "traceDEC":
			return s.traceDEC(APIstub, args)
		case "getDEC":
			return s.getDEC(APIstub, args)
	}
	
	return shim.Error("Invoke Function Not Success.")

}





//creating a digital energy certificate
func (s *SmartContract) createDEC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	// if len(args) != 5 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 5")
	// }

	// "args":["DECID", "BuildingID", "BuildingCategory","FloorArea","HoursOfOccupancy", "EnergyConsumption", "MeterStartDate", "MetereEndDate", "Grade"]


	// ABAC
	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err !=nil {
		return shim.Error("Error retriving user attributes")
	}
	
	if !ok {
		//The client identity does not possess the attributes
		return shim.Error("The client identity does not possess the attributes")
	}

	if val != "buildingowner" && val != "admin" {
		fmt.Println("Attribute role : " + val)
		return shim.Error("Only building owners can create a DEC")
	}


	clientID, _ := cid.GetID(APIstub)
	decID := args[0]
	// buildingID := "util.GenerateUUID()"

	//checking whether the key exists
	decAsBytes, _ := APIstub.GetState(decID)
	if decAsBytes != nil {
		return shim.Error("Key Exist Already")
	}

	// getting the object
	arguments := make([][]byte, 2)
	arguments[0] = []byte("getAsset")
	arguments[1] = []byte(args[1])

	logger.Infof("Getting the identity of the asset")
	response := APIstub.InvokeChaincode("identitycontract", arguments, "mychannel")

	logger.Infof("Received a response from Identity Contract ")
	logger.Infof(fmt.Sprint(response.Status))
	logger.Infof(fmt.Sprint(response.Payload))
	if response.Status != shim.OK || len(response.Payload)==0{
		return shim.Error("Invalid Building ID")
	}

	// logger.Infof(fmt.Sprint(response.Payload))
	// object := Asset{}
	// json.Unmarshal(response.Payload, &object)
	// logger.Infof(object.ID)

	

	//creating the ledger entry
	occupancy, err := strconv.ParseFloat(args[4], 32)
	energy, err :=  strconv.ParseFloat(args[5], 32)
	if err != nil {
		return shim.Error("Invalid Data Types")
	}

	var dec = DEC{DECID: decID, CID: clientID, BuildingID: args[1], Status: "Pending", 
					BuildingCategory : args[2],
					FloorArea: args[3], 
					HoursOfOccupancy: occupancy,
					EnergyConsumption: energy, 
					MeterStartDate: args[6], 
					MetereEndDate: args[7], 
					Grade: args[8]}

	decAsBytes, _ = json.Marshal(dec)
	APIstub.PutState(decID, decAsBytes)

	logger.Infof("Successfully Added")
	return shim.Success(decAsBytes)
}


//updating a digital energy certificate
func (s *SmartContract) updateDEC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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

	if val != "communityverifier" && val != "externalverifier" && val != "admin" {
		fmt.Println("Attribute role : " + val)
		return shim.Error("Insufficient permisions to update the DEC")
	}

	//checking whether the key exists
	decAsBytes, _ := APIstub.GetState(args[0])
	if decAsBytes == nil {
		return shim.Error("Key Doesn't Exist")
	}


	dec := DEC{}

	json.Unmarshal(decAsBytes, &dec)
	dec.Status = args[1]

	decAsBytes, _ = json.Marshal(dec)
	APIstub.PutState(args[0], decAsBytes)

	return shim.Success(decAsBytes)
}


//get DEC
func (s *SmartContract) getDEC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	decAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(decAsBytes)
}



func (t *SmartContract) traceDEC(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	decID := args[0]
	logger.Infof("searching for DEC with id %s", args[0])
	resultsIterator, err := stub.GetHistoryForKey(decID)
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

	logger.Infof("- getHistoryforDEC returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


// // query Access
// func (s *SmartContract) queryAccessRecords(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}
// 	logger.Infof("Quering the ID %s", args[0])
// 	asBytes, _ := APIstub.GetState(args[0])

// 	access := AccessRecord{}
// 	json.Unmarshal(asBytes, &access)
// 	logger.Infof("Sending the object with the this %s", access.Subject)
// 	return shim.Success(asBytes)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}


