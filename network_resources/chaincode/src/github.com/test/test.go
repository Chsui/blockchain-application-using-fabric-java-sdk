package main

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "set" {
		return s.set(stub, args)
	} else if function == "transfer" {
		return s.transfer(stub, args)
	} else if function == "get" {
		return s.get(stub, args)
	}
	return shim.Error("Invalid function name.")
}

func (s *SmartContract) set(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments.")
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("Failed to set asset")
	}
	return shim.Success(nil)
}

func (s *SmartContract) get(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		shim.Error("Error Incorrect arguments.")
	}
	valueAsBytes, _ := stub.GetState(args[0])
	return shim.Success(valueAsBytes)
}

func (s *SmartContract) transfer(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var A, B string
	var Aval, Bval int
	var X int
	var err error
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]
	Avalbytes, err := stub.GetState(A)
	Aval, _ = strconv.Atoi(string(Avalbytes))
	Bvalbytes, err := stub.GetState(B)
	Bval, _ = strconv.Atoi(string(Bvalbytes))
	X, err = strconv.Atoi(args[2])
	if Aval < X {
		return shim.Error("Not enough value")
	}
	Aval = Aval - X
	Bval = Bval + X
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		shim.Error(err.Error())
	}
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract")
	}
}
