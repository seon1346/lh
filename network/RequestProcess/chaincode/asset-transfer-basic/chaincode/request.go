/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type RightProcess struct {
	Key       string     `json:"key"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	Contracts []Contract `json:"Contracts"`
}

type Contract struct {
	Docu_id       uint64 `json:"docu_id"`
	Docu_name     string `json:"docu_name"`
	Document_hash string `json:"document_hash"`
}

//fuction
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, key string, name string, state string) error {

	//marshal
	var request = RightProcess{Key: key, Name: name, State: state}
	requestAsBytes, _ := json.Marshal(request)
	return ctx.GetStub().PutState(key, requestAsBytes)

}

func (s *SmartContract) AddContract(ctx contractapi.TransactionContextInterface, key string, state string, docu_id string, docu_name string, document_hash string) error {
	requestAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return err
	} else if requestAsBytes == nil {
		return fmt.Errorf("User does not exist " + key + "/")
	}

	request := RightProcess{}
	err = json.Unmarshal(requestAsBytes, &request)

	if err != nil {
		return err
	}

	docu_id64, _ := strconv.ParseInt(docu_id, 10, 64)

	Contract := Contract{Docu_id: uint64(docu_id64), Docu_name: docu_name, Document_hash: document_hash}

	request.Contracts = append(request.Contracts, Contract)

	requestAsBytes, err = json.Marshal(request)

	if err != nil {
		return fmt.Errorf("failed to Marshaling:%v", err)

	}

	err = ctx.GetStub().PutState(key, requestAsBytes)

	if err != nil {
		return fmt.Errorf("failed to AddContract %v", err)

	}

	return nil

}

func (s *SmartContract) UpdateState(ctx contractapi.TransactionContextInterface, key string, newstate string) error {

	contractAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return err
	} else if contractAsBytes == nil {
		return fmt.Errorf("User does not exist " + key + "/")
	}

	contract := RightProcess{}
	err = json.Unmarshal(contractAsBytes, &contract)

	if err != nil {
		return err
	}

	contract.State = newstate

	contractAsBytes, err = json.Marshal(contract)

	return ctx.GetStub().PutState(key, contractAsBytes)

	if err != nil {
		return fmt.Errorf("failed to Marshaling:%v", err)

	}

	err = ctx.GetStub().PutState(key, contractAsBytes)

	if err != nil {
		return fmt.Errorf("failed to AddContract %v", err)

	}

	return nil
}
func (s *SmartContract) ReadContract(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	//get value from ctx
	contractAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", fmt.Errorf("failed to read from world state,%s", err.Error())
	}

	if contractAsBytes == nil {
		return "", fmt.Errorf("%s  does not exist", key)

	}

	return string(contractAsBytes[:]), nil

}
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
