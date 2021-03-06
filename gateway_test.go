/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric_connector

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/stretchr/testify/assert"
)

var gw Gateway

func init() {
	var err error
	gw, err = NewGatewayService("./conf/connection-org1.yaml", "User1",
		"./artifacts/crypto-config/peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp", "Org1MSP")
	if err != nil {
		panic(err)
	}
}

func TestGatewayService_CallContract(t *testing.T) {
	payload, err := gw.SubmitTransaction(testChannelId, testCCId, "move", []string{"a", "b", "10", testEventID})
	assert.NoError(t, err)
	t.Log(string(payload))

	result, err := gw.EvaluateTransaction(testChannelId, testCCId, "query", []string{"a"})
	assert.NoError(t, err)
	t.Log(string(result))
}

func TestGatewayService_RegisterContractEvent(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := gw.RegisterChaincodeEvent(ctx, testChannelId, testCCId, testEventID, func(e *fab.CCEvent) {
			t.Logf("ChaincodeEvent receive data: %+v, payload:%s \n", e, string(e.Payload))
		})
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		_, err := gw.SubmitTransaction(testChannelId, testCCId, "move", []string{"a", "b", "10", testEventID})
		assert.NoError(t, err)
		wg.Done()
		time.Sleep(time.Second)
		cancel()
	}()
	wg.Wait()
}

func TestGatewayService_RegisterBlockEvent(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := gw.RegisterBlockEvent(ctx, testChannelId, func(data *TransactionInfo) {
			t.Logf("EventHandler receive data: %+v \n", data)
		})
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		_, err := gw.SubmitTransaction(testChannelId, testCCId, "move", []string{"a", "b", "10", testEventID})
		assert.NoError(t, err)
		wg.Done()
		time.Sleep(time.Second)
		cancel()
	}()
	wg.Wait()
}
