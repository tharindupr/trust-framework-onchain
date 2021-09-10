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
	//"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

type Parameters struct{
	Beta float64 `json:Beta`
	Gamma float64 `json:Gamma`
	AlphaPositive float64 `json:AlphaPositive`
	AlphaNegative float64 `json:AlphaNegative`
}

// type TrustParameters struct{
	
// }

type Prediction struct{
	NodeID string `json:NodeID`
	Timestamp string `json:Timestamp`
	OutPut map[string]bool   `json:"OutPut"`
}

type Reputation struct{
	NodeID string `json:NodeID`
	QoSReputation float64 `json:QoSReputations`
	QoSReputations map[string]float64  `json:QoSReputations`
	Timestamp string `json:Timestamp`
	TrueCount map[string]int  `json:TrueCount`
	FalseCount map[string]int  `json:FalseCount`
	Counter int `json:Counter`
	LatestPrediction map[string]bool   `json:"LatestPrediction"`
}

type Model struct{
	ModelID string `json:ModelID`
	ModelDescription string `json:ModelDescription`
	NodeID string `json:NodeID`
	MalciousPrecision float64 `json:MalciousPrecision`
	MalciousRecall float64 `json:MalciousRecall`
	BenignPrecision float64 `json:BenignPrecision`
	BenignRecall float64 `json:BenignRecall`
	Hash string `json:Hash`
	PositiveTrustScore float64 `json:PositiveTrustScore`
	NegativeTrustScore float64 `json:NegativeTrustScore`
	TrustComposition string `json:TrustComposition`
}

const Gamma = 0.6
const Beta = 2
const AlphaPositive = -1
const AlphaNegative = 1

var logger = flogging.MustGetLogger("models_cc")




// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("Chiancode : datacollection initiated")

	var parameters Parameters
	parameters.Beta = 2
	parameters.AlphaPositive = 1
	parameters.AlphaNegative = -1
	parameters.Gamma = 0.6

	parametersAsBytes, _ := json.Marshal(parameters)
	APIstub.PutState("parameters", parametersAsBytes)
	return shim.Success(nil)
}


// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))
	logger.Errorf("Function name is:  %d", function)
	switch function {
		case "createModel":
			return s.createModel(APIstub, args)
		case "traceModel":
			return s.traceModel(APIstub, args)
		case "getModel":
			return s.getModel(APIstub, args)
		case "reportPrediction":
			return s.reportPrediction(APIstub, args)
		case "trustUpdate":
			return s.trustUpdate(APIstub, args)
			
	}
	return shim.Error("Invoke Function Not Success.")
}


//add a new model 
func (s *SmartContract) createModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var model Model
	json.Unmarshal([]byte(args[0]), &model)


	var PostiveWeight = 2 * (model.BenignPrecision * model.BenignRecall)/(model.BenignPrecision + model.BenignRecall)
	var NegativeWeight = 2 * (model.MalciousPrecision * model.MalciousRecall)/(model.MalciousPrecision + model.MalciousRecall)


	model.PositiveTrustScore = PostiveWeight * AlphaPositive
	model.NegativeTrustScore = NegativeWeight * AlphaNegative

	modeldAsBytes, _ := json.Marshal(model)
	APIstub.PutState(model.ModelID, modeldAsBytes)

	
	logger.Infof("Model Successfully Added")
	return shim.Success(modeldAsBytes)

}

//retrieve a new model 
func (s *SmartContract) getModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tAsBytes)
}

//update Paramters
func (s *SmartContract) updateParameters(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var parameters Parameters
	json.Unmarshal([]byte(args[0]), &parameters)


	parametersAsBytes, _ := json.Marshal(parameters)
	APIstub.PutState("parameters", parametersAsBytes)

	
	logger.Infof("Patameters updated successfully")
	return shim.Success(parametersAsBytes)

}

//add a new detection
func (s *SmartContract) reportPrediction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var prediction Prediction
	json.Unmarshal([]byte(args[0]), &prediction)

	predictionAsBytes, _ := json.Marshal(prediction)
	APIstub.PutState(prediction.NodeID, predictionAsBytes)
	
	logger.Infof("Prediction Successfully Added")
	return shim.Success(predictionAsBytes)

}


func (s *SmartContract) trustUpdate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	logger.Infof("Unmarshalling the payload")
	var prediction Prediction
	json.Unmarshal([]byte(args[0]), &prediction)

	if prediction.OutPut != nil {
		var tempModelID string
		//var M float64
		var MS float64
		var QoS float64

		QoS = 0
		//get the existing reputations of the devices
		reputation := Reputation{}
		reputationAsBytes, _ := APIstub.GetState(prediction.NodeID)
		logger.Infof("Existing Reputation Retrived")
		if reputationAsBytes == nil {
			reputation.QoSReputations = make(map[string]float64)
			reputation.TrueCount = make(map[string]int)
			reputation.FalseCount = make(map[string]int)
			reputation.NodeID = prediction.NodeID
			
		} else{
			json.Unmarshal(reputationAsBytes, &reputation)
		}
		
		//increasing the counter for every trust update
		reputation.Counter++

		//iterating through each model 
		for index, element := range prediction.OutPut {

			logger.Infof("index " + index)
			logger.Infof("element " + fmt.Sprint(element))
			tempModelID = index	
			logger.Infof("Getting the model metadata from Chain")
			modelsAsBytes, _ := APIstub.GetState(tempModelID)
			if modelsAsBytes == nil {
				return shim.Error("Invalid Model ID " + tempModelID)
			}
			model := Model{}
			json.Unmarshal(modelsAsBytes, &model)
			logger.Infof("Unmarshaled the model object")

			//if there's no previous reputation
			
			//trust boostraping
			if reputation.LatestPrediction  == nil{
				logger.Infof("Cold starting the reputation")
				reputation.TrueCount[tempModelID] = 0
				reputation.FalseCount[tempModelID] = 0
				reputation.QoSReputations[tempModelID] = 0.5
				
				if element == true{
					reputation.TrueCount[tempModelID] = 1
				} else{
					reputation.FalseCount[tempModelID] = 1
				}

			} else{
				logger.Infof("Updating the reputation")
				//Increasing the reputation count
				if element == true{
					reputation.TrueCount[tempModelID]++
				} else{
					reputation.FalseCount[tempModelID]++
				}
				
				
				if reputation.Counter == 2{
					logger.Infof("Trust update since the count is 5")
					if reputation.TrueCount[tempModelID] > reputation.FalseCount[tempModelID]{
						MS = Gamma * model.PositiveTrustScore	
						reputation.QoSReputations[tempModelID] = (1-Gamma) * reputation.QoSReputations[tempModelID] + MS

					} else{
						MS = Gamma * model.NegativeTrustScore	
						reputation.QoSReputations[tempModelID] = (1-Gamma) * reputation.QoSReputations[tempModelID] + MS
					}
					
					reputation.TrueCount[tempModelID] = 0
					reputation.FalseCount[tempModelID] = 0
					
				} else{
					logger.Infof("Skipping trust update since the count is not 5")
					reputation.QoSReputations[tempModelID] = reputation.QoSReputations[tempModelID]

				}
				
			}
			logger.Infof("Reputation calculated for model " + tempModelID)
			logger.Infof("Latest Reputation " + fmt.Sprint(reputation.QoSReputations[tempModelID]))
			//Calcualtiong R^q(w) QoS reputation of a device at a given time
			QoS = QoS + reputation.QoSReputations[tempModelID] 
		}

		if reputation.Counter == 2{
			reputation.Counter = 0
		}
		
		logger.Infof("Finished iterations through the models")
		reputation.QoSReputation = QoS/ float64(len(reputation.QoSReputations))
		reputation.Timestamp = prediction.Timestamp
		reputation.LatestPrediction = prediction.OutPut
		reputationAsBytes, _ = json.Marshal(reputation)
		APIstub.PutState(prediction.NodeID, reputationAsBytes)
		logger.Infof("Prediction Successfully Added")
		return shim.Success(reputationAsBytes)
	}

	return shim.Error("Prediction Payload is Null")

}



func (t *SmartContract) traceModel(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]
	logger.Infof("searching for the Perf Contract with id %s", args[0])
	resultsIterator, err := stub.GetHistoryForKey(ID)
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

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}




