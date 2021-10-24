package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"strings"
)

type SmartContract struct {}
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "enrollUser" {
		return s.enrollUser(APIstub, args)
	} else if function == "quitUser" {
		return s.quitUser(APIstub, args)
	} else if function == "getUserInfo" {
		return s.getUserInfo(APIstub, args)
	} else if function == "uploadWeight" {
		return s.uploadWeight(APIstub, args)
	} else if function == "aggregation" {
		return s.aggregation(APIstub, args)
	} else if function == "updateUserToken" {
		return s.updateUserToken(APIstub, args)
	} else if function == "updateGlobalAccuracy" {
		return s.updateGlobalAccuracy(APIstub, args)
	}
	fmt.Println("Please check your function : " + function)
	return shim.Error("Unknown function")
}


func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

var userList []User

type User struct {
	userID string
	userKey string
	Token int
}

// INPUT: user_id, user_private_key | OUTPUT: success or fail
func (s *SmartContract) enrollUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	for _, value := range userList {
		if value.userID == args[0]{
			return shim.Error("userID already exits")
		}
	}

	var user = User{userID: args[0], userKey: args[1], Token: 0}

	// Convert user to []byte
	userJSONBytes, _ := json.Marshal(user)
	err := APIstub.PutState(user.userID, userJSONBytes)
	if err != nil {
		return shim.Error("Failed to enroll user " + user.userID)
	}

	userList = append(userList, user)

	return shim.Success(nil)
}

// INPUT: user_id | OUTPUT: success or fail
func (s *SmartContract) quitUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	input := args[0]

	err := APIstub.DelState(input)
	if err != nil {
		return shim.Error("Failed to delete User")
	}

	for index, value := range userList{
		if value.userID == args[0] {
			userList = append(userList[:index], userList[index+1:]...)
			return shim.Success(nil)
		}
	}

	//return shim.Error("User name ["+args[0]+"] not found!")
	return shim.Success(nil)
}

// INPUT: user_id | OUTPUT: byte(userInfo)
func (s *SmartContract) getUserInfo(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	userAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		fmt.Println(err.Error())
	}

	user := User{}
	json.Unmarshal(userAsBytes, &user)

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	if bArrayMemberAlreadyWritten == true {
		buffer.WriteString(",")
	}
	buffer.WriteString("{\"ID\":")
	buffer.WriteString("\"")
	buffer.WriteString(user.userID)
	buffer.WriteString("\"")

	buffer.WriteString(", \"KEY\":")
	buffer.WriteString("\"")
	buffer.WriteString(user.userKey)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Token\":")
	buffer.WriteString("\"")
	buffer.WriteString(string(user.Token))
	buffer.WriteString("\"")

	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
	buffer.WriteString("]\n")

	return shim.Success(buffer.Bytes())
}

var weightList map[int] []UserWeight

type UserWeight struct {
	uploadUserID string
	weight string
}

// INPUT: user_id, weight, global_round | OUTPUT: success or fail
func (s *SmartContract) uploadWeight(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	inputWeight := UserWeight{uploadUserID: args[0], weight: args[1]}
	inputGlobalRound, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Failed upload Weight")
	}

	weightList[inputGlobalRound] = append(weightList[inputGlobalRound], inputWeight)

	return shim.Success(nil)
}

// INPUT: global_round | OUTPUT: JSON(global_model)
func (s *SmartContract) aggregation(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	currentRound, _ := strconv.Atoi(args[0])

	workerNum := float64(len(weightList))

	totalConv0WeightArray := [...]float64{}
	totalConv0BiasArray := [...]float64{}

	totalConv3WeightArray := [...]float64{}
	totalConv3BiasArray := [...]float64{}

	totalConv6WeightArray := [...]float64{}
	totalConv6BiasArray := [...]float64{}

	totalFc0WeightArray := [...]float64{}
	totalFc0BiasArray := [...]float64{}

	totalFc2WeightArray := [...]float64{}
	totalFc2BiasArray := [...]float64{}

	for i, _ := range weightList {
		conv0weight := jsonToArray(weightList[currentRound][i].weight, "conv0weight")
		conv0bias := jsonToArray(weightList[currentRound][i].weight, "conv0bias")

		conv3weight := jsonToArray(weightList[currentRound][i].weight, "conv3weight")
		conv3bias := jsonToArray(weightList[currentRound][i].weight, "conv3bias")

		conv6weight := jsonToArray(weightList[currentRound][i].weight, "conv6weight")
		conv6bias := jsonToArray(weightList[currentRound][i].weight, "conv6bias")

		fc0weight := jsonToArray(weightList[currentRound][i].weight, "fc0weight")
		fc0bias := jsonToArray(weightList[currentRound][i].weight, "fc0bias")

		fc2weight := jsonToArray(weightList[currentRound][i].weight, "fc2weight")
		fc2bias := jsonToArray(weightList[currentRound][i].weight, "fc2bias")

		for index, _ := range conv0weight {
			totalConv0WeightArray[index] += conv0weight[index]
		}
		for index, _ := range conv0bias {
			totalConv0BiasArray[index] += conv0bias[index]
		}

		for index, _ := range conv3weight {
			totalConv3WeightArray[index] += conv3weight[index]
		}
		for index, _ := range conv3bias {
			totalConv3BiasArray[index] += conv3bias[index]
		}

		for index, _ := range conv6weight {
			totalConv6WeightArray[index] += conv6weight[index]
		}
		for index, _ := range conv6bias {
			totalConv6BiasArray[index] += conv6bias[index]
		}

		for index, _ := range fc0weight {
			totalFc0WeightArray[index] += fc0weight[index]
		}
		for index, _ := range fc0bias {
			totalFc0BiasArray[index] += fc0bias[index]
		}

		for index, _ := range fc2weight {
			totalFc2WeightArray[index] += fc2weight[index]
		}
		for index, _ := range fc2bias {
			totalFc2BiasArray[index] += fc2bias[index]
		}
	}

	for index, _ := range totalConv0WeightArray {
		totalConv0WeightArray[index] = totalConv0WeightArray[index] / workerNum
	}
	for index, _ := range totalConv0BiasArray {
		totalConv0BiasArray[index] = totalConv0BiasArray[index] / workerNum
	}

	for index, _ := range totalConv3WeightArray {
		totalConv3WeightArray[index] = totalConv3WeightArray[index] / workerNum
	}
	for index, _ := range totalConv3BiasArray {
		totalConv3BiasArray[index] = totalConv3BiasArray[index] / workerNum
	}

	for index, _ := range totalConv0WeightArray {
		totalConv0WeightArray[index] = totalConv0WeightArray[index] / workerNum
	}
	for index, _ := range totalConv6BiasArray {
		totalConv6BiasArray[index] = totalConv6BiasArray[index] / workerNum
	}

	for index, _ := range totalFc0WeightArray {
		totalFc0WeightArray[index] = totalFc0WeightArray[index] / workerNum
	}
	for index, _ := range totalFc0BiasArray {
		totalFc0BiasArray[index] = totalFc0BiasArray[index] / workerNum
	}

	for index, _ := range totalFc2WeightArray {
		totalFc2WeightArray[index] = totalFc2WeightArray[index] / workerNum
	}
	for index, _ := range totalFc2BiasArray {
		totalFc2BiasArray[index] = totalFc2BiasArray[index] / workerNum
	}

	globalWeight := make(map[string]interface{})

	globalWeight["conv0weigh"] = totalConv0WeightArray
	globalWeight["conv0bias"] = totalConv0BiasArray

	globalWeight["conv3weigh"] = totalConv3WeightArray
	globalWeight["conv3bias"] = totalConv3BiasArray

	globalWeight["conv6weigh"] = totalConv6WeightArray
	globalWeight["conv6bias"] = totalConv6BiasArray

	globalWeight["fc0weight"] = totalFc0WeightArray
	globalWeight["fc0bias"] = totalFc0BiasArray

	globalWeight["fc2weight"] = totalFc2WeightArray
	globalWeight["fc2bias"] = totalFc2BiasArray

	globalWeightJSON, _ := json.Marshal(globalWeight)

	var buffer bytes.Buffer
	buffer.WriteString(string(globalWeightJSON))

	return shim.Success(buffer.Bytes())
}

func jsonToArray(jsonFile string, layerName string) []float64 {
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonFile), &jsonData)

	stringWeight, _ := jsonData[layerName].(string)
	preprocessingString := strings.Trim(stringWeight, "[]")
	stringArray := strings.Split(preprocessingString, ",")

	var tempArray []float64

	for _, value := range stringArray {
		temp, _ := strconv.ParseFloat(strings.TrimSpace(value),  64)
		tempArray = append(tempArray, temp)
	}

	return tempArray
}

// INPUT: user_id, token | OUTPUT: success or fail
func (s *SmartContract) updateUserToken(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	userID := args[0]
	inputToken,_ := strconv.Atoi(args[1])

	for _, value := range userList {
		if value.userID == userID {
			value.Token = inputToken
		}
	}

	return shim.Success(nil)
}

var globalAccuracy float64

// INPUT: accuracy | OUTPUT: success or fail
func (s *SmartContract) updateGlobalAccuracy(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accuracy,_ := strconv.ParseFloat(args[0], 64)
	globalAccuracy = accuracy

	return shim.Success(nil)
}


