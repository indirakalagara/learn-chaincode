/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
		"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type MSC struct {
	DEDLimit	float64 `json:"dedlimit"`
	OOPLimit	float64	`json:"ooplimit"`
	CFEEDOOP  bool	`json:"cfeedoop"`
	DFEEDOOP  bool	`json:"dfeeoop"`
	INDORFAMIRY	string	`json:"indorfamily"`

}
var msc MSC

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
// Initialize the medical smart contract
		fmt.Println("Initializing Medicat Smart contract")

		var mscData = MSC{DEDLimit: 100, OOPLimit: 150, CFEEDOOP: true, DFEEDOOP: true, INDORFAMIRY: "I"}
		mscDataBytes, err := json.Marshal(&mscData)
		err = stub.PutState("MSCKEY", mscDataBytes)

	  // var mscstr []string
		// mscBytes, _ := json.Marshal(&mscstr)
		//err := stub.PutState("MSCKEY", mscBytes)


		if err != nil {
			fmt.Println("Failed to initialize medical smart contract")
		}

		fmt.Println("Initialization complete")
		return nil, nil

}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}


	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "getmscdata" {
		fmt.Println("invoking getmscdata " + function)
		//msc,err := t.getmscdata(args[0], stub)
		mscBytes,err := t.getmscdata(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the msc")
			return nil, err
		}
		// mscBytes, err1 := json.Marshal(&msc)
		// if err1 != nil {
		// 	fmt.Println("Error marshalling the msc")
		// 	return nil, err1
		// }
		fmt.Println("All success, returning the msc")
		return mscBytes, nil
	}
	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}
						//error

	return nil, errors.New("Received unknown function query: " + function)
}

//func (t *SimpleChaincode) getmscdata(msckey string, stub shim.ChaincodeStubInterface) (MSC, error) {
func (t *SimpleChaincode) getmscdata(msckey string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In getmscdata key is: "+ msckey)

	mscBytes, err := stub.GetState(msckey)
	if err != nil {
		fmt.Println("Error retrieving msc " + msckey)
		return mscBytes, errors.New("Error retrieving msc " + msckey)
	}

	// err = json.Unmarshal(mscBytes, &msc)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling msc " )
	// 	return msc, errors.New("Error unmarshalling msc " )
	// }
	//
	// return msc, nil

	return mscBytes,nil
}
