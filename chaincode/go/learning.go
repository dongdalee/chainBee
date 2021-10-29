package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{}

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
	} else if function == "getUserWeight" {
		return s.getUserWeight(APIstub, args)
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
	userID  string
	userKey string
	Token   string
}

// INPUT: user_id, user_private_key | OUTPUT: success or fail
func (s *SmartContract) enrollUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	for _, value := range userList {
		if value.userID == args[0] {
			return shim.Error("userID already exits")
		}
	}

	var user = User{userID: args[0], userKey: args[1], Token: "0"}

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

	for index, value := range userList {
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

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	inputUserID := args[0]

	var buffer bytes.Buffer

	for i, value := range userList {
		if value.userID == inputUserID {
			buffer.WriteString("[")
			buffer.WriteString("ID:")
			buffer.WriteString(userList[i].userID)
			buffer.WriteString("KEY:")
			buffer.WriteString(userList[i].userKey)
			buffer.WriteString("TOKEN:")
			buffer.WriteString(string(userList[i].Token))
			buffer.WriteString("]")

			return shim.Success(buffer.Bytes())
		}
	}

	return shim.Error("can't search user")
}

var userWeightList map[int][]UserWeight = make(map[int][]UserWeight)

var modelList map[int][]LeNet5 = make(map[int][]LeNet5)

type UserWeight struct {
	uploadUserID string
	weight       LeNet5
}

type LeNet5 struct {
	conv0weight []float64
	conv0bias   []float64
	conv3weight []float64
	conv3bias   []float64
	conv6weight []float64
	conv6bias   []float64
	fc0weight   []float64
	fc0bias     []float64
	fc2weight   []float64
	fc2bias     []float64
}

func (s *SmartContract) uploadWeight(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	inputUserID := args[0]
	jsonWeight := args[1]
	round, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Round number Input Error")
	}

	conv0weightString := jsonToStringList(jsonWeight)["conv0weight"]
	conv0weightStringList := stringToStringList(conv0weightString)
	conv0weightFloat := stringListToFloatList(conv0weightStringList)

	conv0biasString := jsonToStringList(jsonWeight)["conv0bias"]
	conv0biasStringList := stringToStringList(conv0biasString)
	conv0biasFloat := stringListToFloatList(conv0biasStringList)

	conv3weightString := jsonToStringList(jsonWeight)["conv3weight"]
	conv3weightStringList := stringToStringList(conv3weightString)
	conv3weightFloat := stringListToFloatList(conv3weightStringList)

	conv3biasString := jsonToStringList(jsonWeight)["conv3bias"]
	conv3biasStringList := stringToStringList(conv3biasString)
	conv3biasFloat := stringListToFloatList(conv3biasStringList)

	conv6weightString := jsonToStringList(jsonWeight)["conv6weight"]
	conv6weightStringList := stringToStringList(conv6weightString)
	conv6weightFloat := stringListToFloatList(conv6weightStringList)

	conv6biasString := jsonToStringList(jsonWeight)["conv6bias"]
	conv6biasStringList := stringToStringList(conv6biasString)
	conv6biasFloat := stringListToFloatList(conv6biasStringList)

	fc0weightString := jsonToStringList(jsonWeight)["fc0weight"]
	fc0weightStringList := stringToStringList(fc0weightString)
	fc0weighttFloat := stringListToFloatList(fc0weightStringList)

	fc0biasString := jsonToStringList(jsonWeight)["fc0bias"]
	fc0biasStringList := stringToStringList(fc0biasString)
	fc0biasFloat := stringListToFloatList(fc0biasStringList)

	fc2weightString := jsonToStringList(jsonWeight)["fc2weight"]
	fc2weightStringList := stringToStringList(fc2weightString)
	fc2weighttFloat := stringListToFloatList(fc2weightStringList)

	fc2biasString := jsonToStringList(jsonWeight)["fc2bias"]
	fc2biasStringList := stringToStringList(fc2biasString)
	fc2biasFloat := stringListToFloatList(fc2biasStringList)

	//가중치 리스트 생성
	model := LeNet5{conv0weight: conv0weightFloat, conv0bias: conv0biasFloat, conv3weight: conv3weightFloat, conv3bias: conv3biasFloat, conv6weight: conv6weightFloat, conv6bias: conv6biasFloat, fc0weight: fc0weighttFloat, fc0bias: fc0biasFloat, fc2weight: fc2weighttFloat, fc2bias: fc2biasFloat}
	modelList[round] = append(modelList[round], model)

	user := UserWeight{uploadUserID: inputUserID, weight: model}
	userWeightList[round] = append(userWeightList[round], user)

	return shim.Success(nil)
}

func (s *SmartContract) getUserWeight(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	inputUserID := args[0]
	round, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Round number Input Error")
	}

	var result UserWeight

	for _, value := range userWeightList[round] {
		if value.uploadUserID == inputUserID {
			result = value
		}
	}

	conv0weightResult := float64ListToStringList(result.weight.conv0weight)
	conv0biasResult := float64ListToStringList(result.weight.conv0bias)
	conv3weightResult := float64ListToStringList(result.weight.conv3weight)
	conv3biasResult := float64ListToStringList(result.weight.conv3bias)
	conv6weightResult := float64ListToStringList(result.weight.conv6weight)
	conv6biasResult := float64ListToStringList(result.weight.conv6bias)

	fc0weightResult := float64ListToStringList(result.weight.fc0weight)
	fc0biasResult := float64ListToStringList(result.weight.fc0bias)
	fc2weightResult := float64ListToStringList(result.weight.fc2weight)
	fc2biasResult := float64ListToStringList(result.weight.fc2bias)

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("userID:")
	buffer.WriteString(string(result.uploadUserID))
	buffer.WriteString(",")
	buffer.WriteString("weight:")
	buffer.WriteString("[")

	//conv layer
	buffer.WriteString("conv0weight:[")
	for _, value := range conv0weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv0bias:[")
	for _, value := range conv0biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv3weight:[")
	for _, value := range conv3weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv3bias:[")
	for _, value := range conv3biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv6weight:[")
	for _, value := range conv6weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv6bias:[")
	for _, value := range conv6biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")

	//fc layer
	buffer.WriteString("fc0weight:[")
	for _, value := range fc0weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc0bias:[")
	for _, value := range fc0biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc2weight:[")
	for _, value := range fc2weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc2bias:[")
	for _, value := range fc2biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) aggregation(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	round, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Round number Input Error")
	}

	globalModel := LeNet5{
		conv0weight: aggregator(round, "conv0weight"),
		conv0bias:   aggregator(round, "conv0bias"),
		conv3weight: aggregator(round, "conv3weight"),
		conv3bias:   aggregator(round, "conv3bias"),
		conv6weight: aggregator(round, "conv6weight"),
		conv6bias:   aggregator(round, "conv6bias"),
		fc0weight:   aggregator(round, "fc0weight"),
		fc0bias:     aggregator(round, "fc0bias"),
		fc2weight:   aggregator(round, "fc2weight"),
		fc2bias:     aggregator(round, "fc2bias"),
	}

	conv0weightResult := float64ListToStringList(globalModel.conv0weight)
	conv0biasResult := float64ListToStringList(globalModel.conv0bias)
	conv3weightResult := float64ListToStringList(globalModel.conv3weight)
	conv3biasResult := float64ListToStringList(globalModel.conv3bias)
	conv6weightResult := float64ListToStringList(globalModel.conv6weight)
	conv6biasResult := float64ListToStringList(globalModel.conv6bias)
	fc0weightResult := float64ListToStringList(globalModel.fc0weight)
	fc0biasResult := float64ListToStringList(globalModel.fc0bias)
	fc2weightResult := float64ListToStringList(globalModel.fc2weight)
	fc2biasResult := float64ListToStringList(globalModel.fc2bias)

	var buffer bytes.Buffer
	buffer.WriteString("{")

	//conv layer
	buffer.WriteString("conv0weight:[")
	for _, value := range conv0weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv0bias:[")
	for _, value := range conv0biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv3weight:[")
	for _, value := range conv3weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv3bias:[")
	for _, value := range conv3biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv6weight:[")
	for _, value := range conv6weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("conv6bias:[")
	for _, value := range conv6biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")

	//fc layer
	buffer.WriteString("fc0weight:[")
	for _, value := range fc0weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc0bias:[")
	for _, value := range fc0biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc2weight:[")
	for _, value := range fc2weightResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("],")
	buffer.WriteString("fc2bias:[")
	for _, value := range fc2biasResult {
		buffer.WriteString(value)
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())

	return shim.Success(buffer.Bytes())
}

func float64ListToStringList(floatList []float64) []string {
	var result []string

	for _, value := range floatList {
		strValue := strconv.FormatFloat(value, 'f', -1, 64)
		result = append(result, strValue)
	}

	return result
}

func aggregator(round int, layerName string) []float64 {
	var weightAggregation []float64

	if layerName == "conv0weight" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv0weight) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv0weight[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "conv0bias" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv0bias) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv0bias[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "conv3weight" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv3weight) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv3weight[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "conv3bias" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv3bias) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv3bias[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "conv6weight" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv6weight) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv6weight[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "conv6bias" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].conv6bias) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.conv6bias[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "fc0weight" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].fc0weight) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.fc0weight[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "fc0bias" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].fc0bias) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.fc0bias[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "fc2weight" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].fc2weight) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.fc2weight[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else if layerName == "fc2bias" {
		index := 0
		//모델의 길의를 알기 위해 사용
		for index < len(modelList[round][0].fc2bias) {
			var total float64 = 0.0
			for _, model := range modelList[round] {
				weight := model.fc2bias[index]
				total += weight
			}
			weightAggregation = append(weightAggregation, total/float64(len(modelList[round])))
			index += 1
		}
		return weightAggregation
	} else {
		return nil
	}

	return nil
}

//string 가중치 값 -> 스트링값 리스트로 변환
func stringToStringList(data string) []string {
	patten := regexp.MustCompile(`[\[\],]+`)   // 제거할 패턴 선언 | "[ ] ," 제거
	value := patten.ReplaceAllString(data, "") // 문자열 리스트에서 [] 제거
	stringList := strings.Split(value, " ")    // 문자열 리스트 생성

	return stringList
}

//key에 따른 json value 추출
func jsonToStringList(doc string) map[string]string {
	var data map[string]string //json 파일의 가중치값을 문자로 받는다.
	json.Unmarshal([]byte(doc), &data)

	return data
}

// 각 리스트의 스트링값을 float64로 변환
func stringListToFloatList(input []string) []float64 {
	var floatList []float64

	for _, value := range input {
		castValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return nil
		}
		floatList = append(floatList, castValue)
	}

	return floatList
}

// INPUT: user_id, token | OUTPUT: success or fail
func (s *SmartContract) updateUserToken(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	userID := args[0]
	inputToken := args[1]

	for i, value := range userList {
		if value.userID == userID {
			existsToken, _ := strconv.ParseInt(value.Token, 0, 64)
			inputToken, _ := strconv.ParseInt(inputToken, 0, 64)

			currentToken := existsToken + inputToken
			stringCurrentToken := strconv.FormatInt(currentToken, 10)
			userList[i].Token = stringCurrentToken
		}
	}

	return shim.Success(nil)
}

var globalAccuracy map[int64]float64 = make(map[int64]float64)

// INPUT: accuracy | OUTPUT: success or fail
func (s *SmartContract) updateGlobalAccuracy(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	round, _ := strconv.ParseInt(args[0], 0, 64)
	accuracy, _ := strconv.ParseFloat(args[1], 64)

	globalAccuracy[round] = accuracy

	return shim.Success(nil)
}

func (s *SmartContract) getGlobalAccuracy(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	round, _ := strconv.ParseInt(args[0], 0, 64)

	result := strconv.FormatFloat(globalAccuracy[round], 'f', -1, 64)

	var buffer bytes.Buffer
	buffer.WriteString("accuracy: ")
	buffer.WriteString(result)

	return shim.Success(buffer.Bytes())
}
