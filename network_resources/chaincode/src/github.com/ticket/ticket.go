package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Ticket struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Loc      string `json:"loc"`
	Position string `json:"position"`
	Price    string `json:"price"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "createTicket" {
		return s.createTicket(APIstub, args)
	} else if function == "getTicket" {
		return s.getTicket(APIstub, args)
	} else if function == "changeOwner" {
		return s.changeOwner(APIstub, args)
	}
	return shim.Error("Invalid Function name.")
}

func (s *SmartContract) createTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 6 {
		return shim.Error("Error Incorrect arguments.")
	}
	var ticket = Ticket{ID: args[0], Owner: "none", Name: args[1], Date: args[2], Loc: args[3], Position: args[4], Price: args[5]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Error Incorrect arguments.")
	}
	ticketAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(ticketAsBytes)
}

func (s *SmartContract) getAllTicket(APIstub shim.ChaincodeStubInterface) sc.Response {
	keys := APIstub.GetStringArgs()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for index, value := range keys {
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Index\":")
		buffer.WriteString("\"")
		buffer.WriteString(string(index))
		buffer.WriteString("\"")

		buffer.WriteString(", \"Key\":")
		buffer.WriteString(value)
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments")
	}
	ticketAsBytes, _ := APIstub.GetState(args[0])
	ticket := Ticket{}

	json.Unmarshal(ticketAsBytes, &ticket)
	ticket.Owner = args[1]

	ticketAsBytes, _ = json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract")
	}
}
