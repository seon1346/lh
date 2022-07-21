package fabric

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const org1 = "org1"
const org2 = "org2"
const channel_right = "mychannel"
const channel_contract = "mychannel"
const CC = "request"

func GetCC(org, iden_label, channel, CCName string) *gateway.Contract {
	ccpPath := filepath.Join(
		ccproot,
		fmt.Sprintf("connection-%s.yaml", org),
	)
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(gwallet, iden_label),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channel)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract(CCName)

	return contract
}

func InitUser(key, name, status, cert1, cert2 string) error {
	result, err := GetCC(org1, cert1, channel_right, CC).SubmitTransaction("AddUser", key, name, status)
	if err != nil {
		return fmt.Errorf("CC fail[%s]", err.Error())
	}

	if result != nil {
		return fmt.Errorf("AddUser fail[%s]", string(result))
	}
	result, err = GetCC(org1, cert1, channel_contract, CC).SubmitTransaction("AddUser", key, name, status)
	if err != nil {
		return fmt.Errorf("CC fail[%s]", err.Error())
	}
	if result != nil {
		return fmt.Errorf("AddUser fail[%s]", string(result))
	}

	return nil
}

func UpdateRightDoc(key, name, docu_name, docu_id, hash, cert string) error {
	result, err := GetCC(org2, cert, channel_right, CC).SubmitTransaction("AddContract", key, name, docu_id, docu_name, hash)
	if err != nil {
		return fmt.Errorf("CC fail[%s]", err.Error())
	}
	if result != nil {
		return fmt.Errorf("AddUser fail[%s]", string(result))
	}

	return nil
}
func UpdateRightState(key, status, cert string) error {
	result, err := GetCC(org1, cert, channel_right, CC).SubmitTransaction("UpdateState", key, status)
	if err != nil {
		return fmt.Errorf("CC fail[%s]", err.Error())
	}
	if result != nil {
		return fmt.Errorf("UpdateRightState fail[%s]", string(result))
	}

	return nil
}

func GetRight(key, cert string) ([]byte, error) {

	result, err := GetCC(org1, cert, channel_right, CC).EvaluateTransaction("ReadContract", key)
	if err != nil {
		return nil, fmt.Errorf("CC fail[%s]", err.Error())
	}

	return result, nil
}
