package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type SmartContract struct {
}

type Wallet struct {
	Asset string `json:"asset"`
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "setAsset" {
		return s.setAsset(stub, args)
	} else if function == "getAsset" {
		return s.getAsset(stub, args)
	} else if function == "transferAsset" {
		return s.transferAsset(stub, args)
	}
	return shim.Error("Invalid function name.")
}

func (s *SmartContract) setAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments.")
	}
	var wallet = Wallet{Asset: args[1]}

	walletAsBytes, _ := json.Marshal(wallet)
	stub.PutState(args[0], walletAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		shim.Error("Error Incorrect arguments.")
	}
	walletAsBytes, _ := stub.GetState(args[0])
	return shim.Success(walletAsBytes)
}

func (s *SmartContract) transferAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]
	B := args[1]
	X, _ := strconv.Atoi(args[2])

	AWalletAsBytes, _ := stub.GetState(A)
	AWallet := Wallet{}
	json.Unmarshal(AWalletAsBytes, &AWallet)
	assetA, _ := strconv.Atoi(AWallet.Asset)

	BWalletAsBytes, _ := stub.GetState(B)
	BWallet := Wallet{}
	json.Unmarshal(BWalletAsBytes, &BWallet)
	assetB, _ := strconv.Atoi(BWallet.Asset)

	if assetA < X {
		return shim.Error("Not enough Asset")
	}

	assetA = assetA - X
	assetB = assetB + X

	AWallet.Asset = strconv.Itoa(assetA)
	BWallet.Asset = strconv.Itoa(assetB)

	AWalletAsBytes, _ = json.Marshal(AWallet)
	BWalletAsBytes, _ = json.Marshal(BWallet)

	stub.PutState(A, AWalletAsBytes)
	stub.PutState(B, BWalletAsBytes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract")
	}
}
