package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
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

type TicketId struct {
	Num int
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	ticketId := TicketId{Num: 0}
	ticketIdAsBytes, _ := json.Marshal(ticketId)
	APIstub.PutState("lastid", ticketIdAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "createTicket" {
		return s.createTicket(APIstub, args)
	} else if function == "getTicket" {
		return s.getTicket(APIstub, args)
	} else if function == "getAllTickets" {
		return s.getAllTickets(APIstub)
	} else if function == "changeOwner" {
		return s.changeOwner(APIstub, args)
	} else if function == "deleteTicket" {
		return s.deleteTicket(APIstub, args)
	} else if function == "getTransaction" {
		return s.getTransaction(APIstub, args)
	}
	return shim.Error("Invalid Function name.")
}

func setLastId(APIstub shim.ChaincodeStubInterface, Num int) {
	ticketId := TicketId{Num: Num}
	ticketIdAsBytes, _ := json.Marshal(ticketId)
	APIstub.PutState("lastid", ticketIdAsBytes)
}

func getLastId(APIstub shim.ChaincodeStubInterface) int {
	lastIdAsBytes, _ := APIstub.GetState("lastid")
	lastId := TicketId{}
	json.Unmarshal(lastIdAsBytes, &lastId)
	return lastId.Num
}

func (s *SmartContract) createTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Error Incorrect arguments.")
	}
	lastId := getLastId(APIstub)
	lastId += 1
	var ticket = Ticket{ID: strconv.Itoa(lastId), Owner: "none", Name: args[0], Date: args[1], Loc: args[2], Position: args[3], Price: args[4]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(strconv.Itoa(lastId), ticketAsBytes)

	setLastId(APIstub, lastId)

	return shim.Success(nil)
}

func (s *SmartContract) getTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Error Incorrect arguments.")
	}
	ticketAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(ticketAsBytes)
}

func (s *SmartContract) getAllTickets(APIstub shim.ChaincodeStubInterface) sc.Response {
	lastId := getLastId(APIstub)

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for i := 1; i <= lastId; i++ {
		ticketAsBytes, _ := APIstub.GetState(strconv.Itoa(i))
		if string(ticketAsBytes) != "" {
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Id\":")
			buffer.WriteString("\"")
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString("\"")

			buffer.WriteString(", \"Ticket\":")
			buffer.WriteString(string(ticketAsBytes))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
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

func (s *SmartContract) deleteTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Error Incorrect arguments")
	}
	err := APIstub.DelState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) getTransaction(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Error Incorrect arguments")
	}
	cid := APIstub.GetChannelID()
	fmt.Printf(cid)
	/*ledger := peer.GetLedger(cid)
	if ledger == nil {
		return shim.Error("Error GetLedger")
	}
	_, err := ledger.GetTransactionByID(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}*/
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract")
	}
}
