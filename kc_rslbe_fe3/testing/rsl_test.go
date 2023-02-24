package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"git.kaiostech.com/cloud/common/model/oauth2"
	"git.kaiostech.com/cloud/common/model/vibe/vibe_v1_0"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_fe3/testing/test_hawk"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
	"git.kaiostech.com/cloud/thirdparty/github.com/stretchr/testify/assert"
)

const (
	//hosts
	LOCAL_HOST = "127.0.0.1"
	TEST_ENV   = "test.kaiostech.com"
	DEV_ENV    = "dev.kaiostech.com"
	PROD_ENV   = "prod.kaiostech.com"

	AUTH_SERVER_LOCAL = "127.0.0.1"

	//ports
	KC_RSLBE_LOCAL_PORT = "8091"
	AUTH_SERVER_PORT    = "8090"
)

var (
	auth_token string
)

func getAuthToken() (string, error) {

	client := &http.Client{}

	path := fmt.Sprintf("http://%s:%s/v3.0/tokens", AUTH_SERVER_LOCAL, AUTH_SERVER_PORT)
	body := `{
		"grant_type": "password",
		"user_name": "maen.hammour@kaiostech.com",
		"password": "Y+kUoSr6WVx9lnk+HD7rIuh4OEJ6zbiGmckLpwH4GUY=",
		"scope": "core",
		"device": {
			"device_type": 10,
			"brand": "AlcatelOneTouch",
			"model": "GoFlip2",
			"reference": "40441-2AAQUS0",
			"os": "KaiOS",
			"os_version": "2.5",
			"device_id": "123123123144444"
		},
		"application": {
			"id": "kNpFU6NavpPh4e5qnlFz"
		},
		"partner": {
			"id": "FIN_001"
		}
	}`
	body_bytes := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest("POST", path, body_bytes)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	return string(resp_body), nil
}

func TestRegister3I(t *testing.T) {

	var err error
	auth_token, err = getAuthToken()
	if err != nil {
		t.Error(fmt.Sprintf("Error getting auth token from auth service: %s", err))
		return
	}
	partnerID := "PARTNER_002"

	deviceInfo := vibe_v1_0.T_DeviceInfo{
		Imei: "imei_partner_002",
		Imsi: "imsi_partner_002",
		Isdn: "isdn_partner_002",
	}

	path := fmt.Sprintf("http://%s:%s/kc_rsl_be/v1.0/partners/%s/3is", LOCAL_HOST, KC_RSLBE_LOCAL_PORT, partnerID)
	l4g.Debug("Path is: %s", path)
	body, err := json.Marshal(deviceInfo)
	if err != nil {
		t.Error(fmt.Sprintf("Could not marshal request body: %s %s", path, err))
		return
	}

	var token *oauth2.T_AccessToken

	err = json.Unmarshal([]byte(auth_token), &token)
	if err != nil {
		t.Error(fmt.Sprintf("Could not unmarshal auth_token into *oauth2.T_AccessToken: %s", err))
		return
	}

	resp, err := test_hawk.Do(path, "POST", body, token)
	if err != nil {
		t.Error(fmt.Sprintf("Failed to make a request to: %s, %s ", path, err))
		return
	}

	resp_body, _ := ioutil.ReadAll(resp.Body)

	if resp_body == nil {
		t.Error(fmt.Sprintf("Null body received for: %s ", path))
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error(fmt.Sprintf("Received status code: %v Response body: %s", resp.StatusCode, string(resp_body)))
		return
	}

	assert.Equal(t, resp.StatusCode, http.StatusCreated, fmt.Sprintf("Invalid status code received %v when calling: %s", resp.StatusCode, path))

	l4g.Debug("Received body: %s", resp_body)

	received_device_info, err := vibe_v1_0.JsonToDeviceInfo(resp_body)
	if err != nil {
		t.Error(fmt.Sprintf("Error unmarshalling received body into DeviceInfo struct: %s, %s ", err, resp_body))
		return
	}

	fmt.Printf("Received device info: %s", received_device_info)

	if received_device_info == nil {
		t.Error(fmt.Sprintf("Could not unmarshal received_device_info: %s ", err))
		return
	}
	assert.Equal(t, received_device_info.Imei, deviceInfo.Imei, "The two IMEIs should be the same")
	assert.Equal(t, received_device_info.Imsi, deviceInfo.Imsi, "The two IMSIs should be the same")
	assert.Equal(t, received_device_info.Isdn, deviceInfo.Isdn, "The two ISDNs should be the same")

}
func TestRslCommand(t *testing.T) {

}
func TestUpdate2I(t *testing.T) {
	var err error
	auth_token, err = getAuthToken()
	if err != nil {
		t.Error(fmt.Sprintf("Error getting auth token from auth service: %s", err))
		return
	}
	imei := "some_imei"

	deviceInfo := vibe_v1_0.T_DeviceInfo{
		Imsi: "updated_imei_partner_001",
		Isdn: "updated_isdn_partner_001",
	}

	path := fmt.Sprintf("http://%s:%s/kc_rsl_be/v1.0/devices/update/%s", LOCAL_HOST, KC_RSLBE_LOCAL_PORT, imei)
	l4g.Debug("Path is: %s", path)
	body, err := json.Marshal(deviceInfo)
	if err != nil {
		t.Error(fmt.Sprintf("Could not marshal request body: %s %s", path, err))
		return
	}

	var token *oauth2.T_AccessToken

	err = json.Unmarshal([]byte(auth_token), &token)
	if err != nil {
		t.Error(fmt.Sprintf("Could not unmarshal auth_token into *oauth2.T_AccessToken: %s", err))
		return
	}

	resp, err := test_hawk.Do(path, "PUT", body, token)
	if err != nil {
		t.Error(fmt.Sprintf("Failed to make a request to: %s, %s ", path, err))
		return
	}

	resp_body, _ := ioutil.ReadAll(resp.Body)

	if resp_body == nil {
		t.Error(fmt.Sprintf("Null body received for: %s ", path))
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error(fmt.Sprintf("Received status code: %v Response body: %s", resp.StatusCode, string(resp_body)))
		return
	}

	assert.Equal(t, resp.StatusCode, http.StatusCreated, fmt.Sprintf("Invalid status code received %v when calling: %s", resp.StatusCode, path))

	l4g.Debug("Received body: %s", resp_body)

	received_device_info, err := vibe_v1_0.JsonToDeviceInfo(resp_body)
	if err != nil {
		t.Error(fmt.Sprintf("Error unmarshalling received body into DeviceInfo struct: %s, %s ", err, resp_body))
		return
	}

	fmt.Printf("Received device info: %s", received_device_info)

	if received_device_info == nil {
		t.Error(fmt.Sprintf("Could not unmarshal received_device_info: %s ", err))
		return
	}
	assert.Equal(t, received_device_info.Imei, imei, "The two IMEIs should be the same")
	assert.Equal(t, received_device_info.Imsi, deviceInfo.Imsi, "The two IMSIs should be the same")
	assert.Equal(t, received_device_info.Isdn, deviceInfo.Isdn, "The two ISDNs should be the same")
}
func TestTransferOwnership(t *testing.T) {
	var err error
	auth_token, err = getAuthToken()
	if err != nil {
		t.Error(fmt.Sprintf("Error getting auth token from auth service: %s", err))
		return
	}
	imei := "imei_partner_001"
	pid := "PARTNER_002"

	path := fmt.Sprintf("http://%s:%s/kc_rsl_be/v1.0/devices/transfer_ownership/%s/%s", LOCAL_HOST, KC_RSLBE_LOCAL_PORT, imei, pid)
	l4g.Debug("Path is: %s", path)

	if err != nil {
		t.Error(fmt.Sprintf("Could not marshal request body: %s %s", path, err))
		return
	}

	var token *oauth2.T_AccessToken

	err = json.Unmarshal([]byte(auth_token), &token)
	if err != nil {
		t.Error(fmt.Sprintf("Could not unmarshal auth_token into *oauth2.T_AccessToken: %s", err))
		return
	}

	resp, err := test_hawk.Do(path, "POST", nil, token)
	if err != nil {
		t.Error(fmt.Sprintf("Failed to make a request to: %s, %s ", path, err))
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error(fmt.Sprintf("Received status code: %v ", resp.StatusCode))
		return
	}

	assert.Equal(t, resp.StatusCode, http.StatusCreated, fmt.Sprintf("Invalid status code received %v when calling: %s", resp.StatusCode, path))
}
func TestTransferState(t *testing.T) {
	var err error
	auth_token, err = getAuthToken()
	if err != nil {
		t.Error(fmt.Sprintf("Error getting auth token from auth service: %s", err))
		return
	}
	imei_1 := "imei_1"
	imei_2 := "imei_2"

	path := fmt.Sprintf("http://%s:%s/kc_rsl_be/v1.0/devices/transfer_state/%s/%s", LOCAL_HOST, KC_RSLBE_LOCAL_PORT, imei_1, imei_2)
	l4g.Debug("Path is: %s", path)

	if err != nil {
		t.Error(fmt.Sprintf("Could not marshal request body: %s %s", path, err))
		return
	}

	var token *oauth2.T_AccessToken

	err = json.Unmarshal([]byte(auth_token), &token)
	if err != nil {
		t.Error(fmt.Sprintf("Could not unmarshal auth_token into *oauth2.T_AccessToken: %s", err))
		return
	}

	resp, err := test_hawk.Do(path, "POST", nil, token)
	if err != nil {
		t.Error(fmt.Sprintf("Failed to make a request to: %s, %s ", path, err))
		return
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error(fmt.Sprintf("Received status code: %v ", resp.StatusCode))
		return
	}

	assert.Equal(t, resp.StatusCode, http.StatusCreated, fmt.Sprintf("Invalid status code received %v when calling: %s", resp.StatusCode, path))
}
