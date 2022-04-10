package main

import (
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

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createTicket" {
		return s.createTicket(stub, args)
	} else if function == "getTicket" {
		return s.getTicket(stub, args)
	} else if function == "changeOwner" {
		return s.changeOwner(stub, args)
	}
	return shim.Error("Invalid Function name.")
}

func (s *SmartContract) createTicket(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 6 {
		return shim.Error("Error Incorrect arguments.")
	}
	var ticket = Ticket{ID: args[0], Owner: "none", Name: args[1], Date: args[2], Loc: args[3], Position: args[4], Price: args[5]}

	ticketAsBytes, _ := json.Marshal(ticket)
	stub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getTicket(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Error Incorrect arguments.")
	}
	ticketAsBytes, _ := stub.GetState(args[0])
	return shim.Success(ticketAsBytes)
}

func (s *SmartContract) changeOwner(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments")
	}
	ticketAsBytes, _ := stub.GetState(args[0])
	ticket := Ticket{}

	json.Unmarshal(ticketAsBytes, &ticket)
	ticket.Owner = args[1]

	ticketAsBytes, _ = json.Marshal(ticket)
	stub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract")
	}
}
