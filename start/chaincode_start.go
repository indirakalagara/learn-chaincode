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

// InsuranceChaincode example simple Chaincode implementation
type InsuranceChaincode struct {
}

type MSC struct {
	DEDLimit	float64 `json:"dedlimit"`
	OOPLimit	float64	`json:"ooplimit"`
	CFEEDOOP  bool	`json:"cfeedoop"`
	DFEEDOOP  bool	`json:"dfeeoop"`
	INDORFAMIRY	string	`json:"indorfamily"`

}
type AccumShare struct {
	Claims struct {
		PolicyID string `json:"PolicyID"`
		SubscriberID string `json:"SubscriberID"`
		PolicyStartDate string `json:"PolicyStartDate"`
		PolicyEndDate string `json:"PolicyEndDate"`
		PolicyType string `json:"PolicyType"`
		DeductibleBalance string `json:"DeductibleBalance"`
		OOPBalance string `json:"OOPBalance"`
		Claim struct {
			ClaimID string `json:"ClaimID"`
			MemberID string `json:"MemberID"`
			CreateDTTM string `json:"CreateDTTM"`
			LastUpdateDTTM string `json:"LastUpdateDTTM"`
			Transaction struct {
				TransactionID string `json:"TransactionID"`
				Accumulator struct {
					Type string `json:"Type"`
					Amount string `json:"Amount"`
					UoM string `json:"UoM"`
				} `json:"Accumulator"`
				Participant string `json:"Participant"`
				TotalTransactionAmount string `json:"TotalTransactionAmount"`
				UoM string `json:"UoM"`
			} `json:"Transaction"`
			TotalClaimAmount string `json:"TotalClaimAmount"`
			UoM string `json:"UoM"`
		} `json:"Claim"`
	} `json:"Claims"`
}
//contractstruct  - data struct

var msc MSC

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *InsuranceChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the medical smart contract
	fmt.Println("Initializing  Smart contract")

	var mscData = MSC{DEDLimit: 100, OOPLimit: 150, CFEEDOOP: true, DFEEDOOP: true, INDORFAMIRY: "I"}
	mscDataBytes, err := json.Marshal(&mscData)
	err = stub.PutState("MSCKEY", mscDataBytes)

	//Initialize AccumShare
	accumShareJson := `{   "Claims": {      "PolicyID": "1266363",      "SubscriberID": "10003",      "PolicyStartDate": "05-Jan-2016",      "PolicyEndDate": "31-Dec-2017",      "PolicyType": "Individual",            "DeductibleBalance":"600",      "OOPBalance":"300",    	  "BalanceUoM":"Dollars", 	 	        "Claim": {         "ClaimID": "18738936",         "MemberID": "10003",         "CreateDTTM": "11-Jan-2017",         "LastUpdateDTTM": "11-Jan-2017",         "Transaction": {            "TransactionID": "36563856",            "Accumulator": {               "Type": "Deductible",                              "Amount": "200",               "UoM": "Dollars"            },            "Participant": "Medical",            "TotalTransactionAmount": "200",            "UoM": "Dollars"         },         "TotalClaimAmount": "200",         "UoM": "Dollars"      }   }}`
	var accumShare AccumShare
	err = json.Unmarshal([]byte(accumShareJson), accumShare)
	if err != nil {
		fmt.Println("Failed to Unmarshal  Accumshare ")
	}
	err = stub.PutState("10003", []byte(accumShareJson))
	if err != nil {
		fmt.Println("Failed to initialize  smart contract")
	}

	fmt.Println("Initialization complete")
	return nil, nil

}

// Invoke is our entry point to invoke a chaincode function
func (t *InsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *InsuranceChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "getMscData" {
		fmt.Println("invoking getMscData " + function)
		//msc,err := t.getMscData(args[0], stub)
		mscBytes,err := t.getMscData(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the msc")
			return nil, err
		}

		fmt.Println("All success, returning the msc")
		return mscBytes, nil
	}
	if function == "getAccumShare" {
		fmt.Println("invoking getAccumShare " + function)
		//msc,err := t.getMscData(args[0], stub)
		accumShareBytes,err := t.getAccumShare(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the AccumShare")
			return nil, err
		}

		fmt.Println("All success, returning the accumShare")
		return accumShareBytes, nil
	}

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}

	return nil, errors.New("Received unknown function query: " + function)
}

//func (t *InsuranceChaincode) getmscdata(msckey string, stub shim.ChaincodeStubInterface) (MSC, error) {
func (t *InsuranceChaincode) getMscData(msckey string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In getMscData key is: "+ msckey)

	mscBytes, err := stub.GetState(msckey)
	if err != nil {
		fmt.Println("Error retrieving msc " + msckey)
		return mscBytes, errors.New("Error retrieving msc " + msckey)
	}

	return mscBytes,nil
}
func (t *InsuranceChaincode) getAccumShare(subscriberID string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In getAccumShare subscriberID is: "+ subscriberID)

	accumShareBytes, err := stub.GetState(subscriberID)
	if err != nil {
		fmt.Println("Error retrieving AccumShare " + subscriberID)
		return nil, errors.New("Error retrieving AccumShare " + subscriberID)
	}

	var accumShare AccumShare
	err = json.Unmarshal(accumShareBytes, &accumShare)
	//accumShareJson, err := json.Marshal(accumShare)
	//fmt.Println("accumSharejson  is : " , accumShareJson);

	fmt.Println("accumShare  is : " , accumShare);
	fmt.Println("accumShare deductible balance is : " , accumShare.Claims.DeductibleBalance);

	return accumShareBytes,nil
}

func (t *InsuranceChaincode) setMscData(msckey string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In setmscdata for key is: "+ msckey)

	var mscData = MSC{DEDLimit: 100, OOPLimit: 150, CFEEDOOP: true, DFEEDOOP: true, INDORFAMIRY: "I"}
	mscDataBytes, err := json.Marshal(&mscData)
	err = stub.PutState("MSCKEY", mscDataBytes)

	if err != nil {
		fmt.Println("Failed to add  medical smart contract")
	}

	return nil,nil
}
