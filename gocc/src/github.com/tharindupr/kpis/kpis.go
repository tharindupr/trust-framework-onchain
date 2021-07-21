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




type Telemetry struct{
	Unit string `json:"Unit"`
	Week1 map[string]float64  `json:"Week1"`
	Week2 map[string]float64  `json:"Week2"`
	Week3 map[string]float64  `json:"Week3"`
	Week4 map[string]float64  `json:"Week4"`
	Week5 map[string]float64  `json:"Week5"`
}

type EnergyRecord struct{
	TimeStamp int `json:"TimeStamp"`
	Year int `json:"Year"`
	Month int `json:"Month"`
	WeeksPerMonth int `json:"WeeksPerMonth"`
	BuildingID  string `json:"BuildingID"`
	Readings  Telemetry `json:"Readings"`
	Status bool `json:"Status"`
}


type Perfomance struct{
	BuildingID  string `json:"BuildingID"`
	TotalTargetUsage  float64 `json:"TotalTargetUsage"`
	IndividualTargetUsage Telemetry `json:"IndividualTargetUsage"`
	Baseline float64 `json:"Baseline"`
	IndividualBaseline Telemetry `json:"IndividualBaseline"`
}

type KPI struct{
	BuildingID  string `json:"BuildingID"`
	Year int `json:"Year"`
	Month int `json:"Month"`
	TotalUsage float64 `json:"TotalUsage"`
	TotalReduction float64 `json:"Reduction"`
	WeeklyReductions Telemetry `json:"WeeklyReductions"`
	TotalTargetUsage float64 `json:"TotalTargetUsage"`
	Baseline float64 `json:"Baseline"`
	IndividualTargetUsage Telemetry `json:"IndividualTargetUsage"`
	IndividualBaseline Telemetry `json:"IndividualBaseline"`
}


var logger = flogging.MustGetLogger("kpis_cc")


// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("Chiancode : datacollection initiated")
	return shim.Success(nil)
}


// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))
	logger.Errorf("Function name is:  %d", function)
	switch function {
		case "createEnergyRecord":
			return s.createEnergyRecord(APIstub, args)
		case "addWeeklyEnergyData":
			return s.addWeeklyEnergyData(APIstub, args)
		case "createPerformanceContract":
			return s.createPerformanceContract(APIstub, args)
		case "getRecordByKey":
			return s.getRecordByKey(APIstub, args)
		case "traceTransactionHistory":
			return s.traceTransactionHistory(APIstub, args)			
	}
	return shim.Error("Invoke Function Not Success.")
}





//creating a digital energy certificate
func (s *SmartContract) createEnergyRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var data EnergyRecord
	json.Unmarshal([]byte(args[0]), &data)
	
	// getting the object
	arguments := make([][]byte, 2)
	arguments[0] = []byte("getAsset")
	arguments[1] = []byte(data.BuildingID)

	logger.Infof("Getting the identity of the asset")
	response := APIstub.InvokeChaincode("identitycontract", arguments, "mychannel")

	logger.Infof("Received a response from Identity Contract ")
	logger.Infof(fmt.Sprint(response.Status))
	logger.Infof(fmt.Sprint(response.Payload))
	if response.Status != shim.OK || len(response.Payload)==0{
		return shim.Error("Invalid Building ID")
	}


	var index = data.BuildingID + strconv.Itoa(data.Year) + strconv.Itoa(data.Month)

	energyRecordAsBytes, _ := json.Marshal(data)
	APIstub.PutState(index, energyRecordAsBytes)

	//calcualte energy KPIs
	if data.Status == true{
		calculateKPIs(APIstub, data)
	}
	
	logger.Infof("Successfully Added")
	return shim.Success(energyRecordAsBytes)
}


//Adding weekly data
func (s *SmartContract) addWeeklyEnergyData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//checking whether the key exists
	energyRecordAsBytes, _ := APIstub.GetState(args[0])
	if energyRecordAsBytes == nil {
		return shim.Error("Key Doesn't Exist")
	}


	energyrecord := EnergyRecord{}
	json.Unmarshal(energyRecordAsBytes, &energyrecord)

	var weeklydata = map[string]float64{};
	json.Unmarshal([]byte(args[2]), &weeklydata)

	if args[1] == "Week1" {
		energyrecord.Readings.Week1 = weeklydata
	} else if args[1] == "Week2" {
		energyrecord.Readings.Week2 = weeklydata
	} else if args[1] == "Week3" {
		energyrecord.Readings.Week3 = weeklydata
	} else if args[1] == "Week4" {
		energyrecord.Readings.Week4 = weeklydata
	} else if args[1] == "Week5" {
		energyrecord.Readings.Week5 = weeklydata
	} else {
		return shim.Error("Invalid Week")
	}

	var index = energyrecord.BuildingID + strconv.Itoa(energyrecord.Year) + strconv.Itoa(energyrecord.Month)

	//check whether month energy records are completed
	if energyrecord.WeeksPerMonth == 4 && energyrecord.Readings.Week1 != nil && energyrecord.Readings.Week2 != nil && energyrecord.Readings.Week3 != nil && energyrecord.Readings.Week4 != nil{
		
		energyrecord.Status = true
	}

	if energyrecord.Readings.Week1 != nil && energyrecord.Readings.Week2 != nil && energyrecord.Readings.Week3 != nil && energyrecord.Readings.Week4 != nil && energyrecord.Readings.Week5 != nil{
		
		energyrecord.Status = true
	}


	//calcualte energy KPIs
	if energyrecord.Status == true{

		calculateKPIs(APIstub, energyrecord)

	}

	energyRecordAsBytes, _ = json.Marshal(energyrecord)
	APIstub.PutState(index, energyRecordAsBytes)

	
	return shim.Success(energyRecordAsBytes)
}


func calculateKPIs(APIstub shim.ChaincodeStubInterface, energyrecord EnergyRecord){

	var totalUsageByEnergy map[string]float64 = calculateMonthlyUsage(energyrecord)
	var index = energyrecord.BuildingID + strconv.Itoa(energyrecord.Year) + strconv.Itoa(energyrecord.Month)
	logger.Infof(fmt.Sprint(totalUsageByEnergy))

	var totalUsage float64 = 0
	for _, value := range totalUsageByEnergy {

		totalUsage = totalUsage + value
	}

	var kpi KPI
	//getting the performance contract
	logger.Infof("Getting the performance contract of the ID " + energyrecord.BuildingID)
	contractAsBytes, _ := APIstub.GetState(energyrecord.BuildingID)
	logger.Infof(fmt.Sprint(contractAsBytes))
	if contractAsBytes == nil {
		// if contract not available
		var telem Telemetry
		kpi.TotalUsage = totalUsage
		kpi.TotalReduction = 0
		kpi.TotalTargetUsage = 0
		kpi.WeeklyReductions = telem
		kpi.Baseline = 0
	
	} else{
		contract := Perfomance{}
		json.Unmarshal(contractAsBytes, &contract)
		kpi.TotalUsage = totalUsage
		kpi.TotalReduction = (contract.TotalTargetUsage - totalUsage)/contract.TotalTargetUsage * 100
		kpi.TotalTargetUsage = contract.TotalTargetUsage 
		kpi.Baseline = contract.Baseline 
		kpi.IndividualTargetUsage = contract.IndividualTargetUsage
		kpi.IndividualBaseline = contract.IndividualBaseline
		kpi.WeeklyReductions = calculateWeeklyReduction(energyrecord.Readings, contract.IndividualTargetUsage)

	}
	
	kpi.BuildingID = energyrecord.BuildingID
	kpi.Year = energyrecord.Year
	kpi.Month = energyrecord.Month

	

	kpiAsBytes, _ := json.Marshal(kpi)
	APIstub.PutState("kpi_"+index, kpiAsBytes)

}


//creating a digital energy certificate
func (s *SmartContract) createPerformanceContract(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	//creating an object from the attribute array
	//var jsonObj interface{}
	var perf Perfomance
	json.Unmarshal([]byte(args[0]), &perf)



	performanceAsBytes, _ := json.Marshal(perf)
	APIstub.PutState(perf.BuildingID, performanceAsBytes)

	logger.Infof("Successfully Added")
	return shim.Success(performanceAsBytes)
}



//getRecordByKey
func (s *SmartContract) getRecordByKey(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tAsBytes)
}





func calculateMonthlyUsage(energyrecord EnergyRecord) map[string]float64 {

	logger.Infof("Inside the calculateMonthlyUsage Function")

	var weeklydata = [5] map[string]float64 {energyrecord.Readings.Week1, energyrecord.Readings.Week2, energyrecord.Readings.Week3, energyrecord.Readings.Week4, energyrecord.Readings.Week5}
	energy := make(map[string]float64)

	logger.Infof("Begin Iteration")
	logger.Infof(fmt.Sprint(weeklydata[0]))

	for _, week := range weeklydata {
		if week != nil{
			for key, value := range week {

				if _, ok := energy[key]; ok {
					energy[key] = energy[key] + value
				} else {
					energy[key] = value
				}

			}
		}
	}
	logger.Infof("Finished Iterations")

	return energy
}


func calculateWeeklyReduction(readings Telemetry, benchmark Telemetry) Telemetry{

	logger.Infof("Calculating weekly reductions")

	var weeklyreductions [5] map[string]float64 
	var weeklydata = [5] map[string]float64 {readings.Week1, readings.Week2, readings.Week3, readings.Week4, readings.Week5}
	var weeklybenchmark = [5] map[string]float64 {benchmark.Week1, benchmark.Week2, benchmark.Week3, benchmark.Week4, benchmark.Week5}

	for index, week := range weeklydata {
		temp := make(map[string]float64)
		if week != nil && weeklybenchmark[index] !=nil{
			for key, value := range week {

				temp[key] = (weeklybenchmark[index][key]-value)/weeklybenchmark[index][key] * 100

			}
			weeklyreductions[index] = temp
		} else{
			weeklyreductions[index] = nil
		}
	}

	var weeklyReductionOutput Telemetry
	weeklyReductionOutput.Unit = "kwh"
	weeklyReductionOutput.Week1 = weeklyreductions[0]
	weeklyReductionOutput.Week2 = weeklyreductions[1]
	weeklyReductionOutput.Week3 = weeklyreductions[2]
	weeklyReductionOutput.Week4 = weeklyreductions[3]
	weeklyReductionOutput.Week5 = weeklyreductions[4]

	return weeklyReductionOutput
}

func (t *SmartContract) traceTransactionHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {

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


