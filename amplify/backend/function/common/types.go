package common

import "fmt"

type ManagerInfo struct {
	DeviceManagerID int64  `json:"devicemanagerid"`
	ManagerHostName string `json:"managerhostname"`
	SerialNumber    string `json:"serialnumber"`
	DeviceID        int64  `json:"deviceid"`
}

type DbCreds struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	DbInstanceIdentifier string `json:"dbinstanceidentifier"`
}

func (creds *DbCreds) GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", creds.Host, creds.Port, creds.Username, creds.Password, creds.DbInstanceIdentifier)
}

func (creds *DbCreds) GetEngine() string {
	return creds.Engine
}

type ManagerInfoResponse struct {
	StatusCode  int         `json:"statuscode"`
	ManagerInfo ManagerInfo `json:"managerinfo"`
	Body        RespBody    `json:"body"`
}

type RespBody struct {
	Result  string `json:"result"`
	Payload []byte `json:"payload"`
}

type Parameters struct {
	SerialNumber string  `json:"serialnumber"`
	Command      Command `json:"command"`
}

type Command struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Value string `json:"value"`
}
