package fabric

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var gwallet *gateway.Wallet
var ccproot string = "../ccp/"

func populateWallet(wallet *gateway.Wallet, user, org, mspid string) error {
	credPath := filepath.Join(

		"../",
		"msp",
		fmt.Sprintf("%s.example.com", org),
		"users",
	)

	certPath := filepath.Join(credPath, "signcerts", fmt.Sprintf("%s@%s.example.com-cert.pem", user, org))
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(fmt.Sprintf("../msp/%s.example.com/users/keystore/", org))
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspid, string(cert), string(key))

	err = wallet.Put(fmt.Sprintf("%s_%s", user, org), identity)
	if err != nil {
		return err
	}
	return nil
}
func Init_Wallet() {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}
	gwallet = wallet
	err = populateWallet(wallet, "User1", "org2", "Org2MSP")
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	err = populateWallet(wallet, "User1", "org1", "Org1MSP")
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
}
func ConTest() {

	// Path to the network config (CCP) file
	ccpPath := filepath.Join(
		ccproot,
		"connection-org2.yaml",
	)

	// Connect to the gateway peer(s) using the network config and identity in the wallet
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(gwallet, "User1_org2"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()
	fmt.Print(gw)
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("request")

	result, err := contract.EvaluateTransaction("ReadContract", "1023")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
	/*
		result, err = contract.SubmitTransaction("createCar", "CAR10", "VW", "Polo", "Grey", "Mary")
		if err != nil {
			fmt.Printf("Failed to submit transaction: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(result))

		result, err = contract.EvaluateTransaction("queryCar", "CAR10")
		if err != nil {
			fmt.Printf("Failed to evaluate transaction: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(result))

		_, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
		if err != nil {
			fmt.Printf("Failed to submit transaction: %s\n", err)
			os.Exit(1)
		}

		result, err = contract.EvaluateTransaction("queryCar", "CAR10")
		if err != nil {
			fmt.Printf("Failed to evaluate transaction: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(result))
	*/

}
