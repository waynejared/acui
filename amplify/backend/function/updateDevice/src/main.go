package main

/*to successfully deploy, in the terminal set
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

Then use this build command
go build -o main

You should have a file main with no .exe
zip that up and upload it to your lambda
make sure the Handler varialbe for the lambda is set to 'main'
*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"log"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"github.com/waynejared/eac_ui/amplify/backend/function/common"
)

func HandleRequest(ctx context.Context, manager common.ManagerInfo) (string, error) {
	var statuscode int
	var result string
	//	var tinyData TinyDeviceData
	var respStruct common.ManagerInfoResponse
	var queryArgs []any

	log.Println("starting handler")
	log.Println(manager)
	sqlStmt := `select json_build_object('managerhostname', address, 'deviceid', devicemanagerdevices.id) as manager from eac.devicemanager join eac.devicemanagerdevices on devicemanager.id = devicemanagerid where serialnumber = ($1);`
	queryArgs = append(queryArgs, manager.SerialNumber)
	common.ExecuteSQL(sqlStmt, queryArgs, &respStruct)

	if respStruct.StatusCode == 200 {
		json.Unmarshal(respStruct.Body.Payload, &manager)

		//	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=enable", host, port, user, password, dbname)
		managerURL := "http://" + manager.ManagerHostName + ":8080/config/updatedevice"
		log.Println("Starting deviceUpdate handler for:" + manager.SerialNumber + " at: " + managerURL)

		sqlStmt := "select * from eac.getdeviceconfig($1);"
		common.ExecuteSQL(sqlStmt, queryArgs, &respStruct)
		log.Println(string(respStruct.Body.Payload))
		requestBody := bytes.NewReader(respStruct.Body.Payload)
		request, err := http.NewRequest("POST", managerURL, requestBody)
		if err != nil {
			result = err.Error()
			fmt.Println(result)
			statuscode = 500
		} else {
			request.Header.Add("Content-type", "application/json")
			client := &http.Client{}
			fmt.Println("about to make request")
			response, err := client.Do(request)
			fmt.Println("finished request")
			if err != nil {
				result = err.Error()
				fmt.Println(result)
				statuscode = 400
			} else {
				defer response.Body.Close()
				fmt.Println(result)
				result = "success"
				statuscode = 200
			}
		}
	}
	returnResult := fmt.Sprintf(`{"statuscode":%d,"body":{"result":"%s"}}`, statuscode, result)
	log.Println(returnResult)

	return returnResult, nil
}

func main() {
	lambda.Start(HandleRequest)
}
