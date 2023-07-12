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
	"io"
	"net/http"

	"log"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"github.com/waynejared/eac_ui/amplify/backend/function/common"
)

func HandleRequest(ctx context.Context, inbound map[string]interface{}) (string, error) {
	var parameters common.Parameters
	var statuscode int
	var result string
	var url string
	var queryArgs []any
	var respStruct common.ManagerInfoResponse
	var manager common.ManagerInfo

	log.Println("starting requestAccess handler")
	log.Println(inbound)
	tempJson, _ := json.Marshal(inbound["body"])
	log.Println(string(tempJson))
	json.Unmarshal(tempJson, &result)
	log.Println(result)
	json.Unmarshal([]byte(result), &parameters)
	log.Println(parameters)
	log.Println("Done Printing")
	//	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=enable", host, port, user, password, dbname)
	sqlStmt := `select json_build_object('managerhostname', address) as manager from eac.devicemanager join eac.devicemanagerdevices on devicemanager.id = devicemanagerid where serialnumber = ($1);`

	queryArgs = append(queryArgs, parameters.SerialNumber)

	common.ExecuteSQL(sqlStmt, queryArgs, &respStruct)

	if respStruct.StatusCode == 200 {
		json.Unmarshal(respStruct.Body.Payload, &manager)

		url = "http://" + manager.ManagerHostName + ":8080" + parameters.Command.Path
		log.Println("managerURL: " + url)
		body := []byte(`{"serialnumber":"` + parameters.SerialNumber + `","` + parameters.Command.Name + `":"` + parameters.Command.Value + `"}`)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			result = err.Error()
			statuscode = 500
		} else {
			request.Header.Add("Content-type", "application/json")
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				result = err.Error()
				statuscode = 400
			} else {
				bodyBytes, err := io.ReadAll(response.Body)
				if err != nil {
					log.Println("ReadAll of response failed with " + err.Error())
				} else {
					defer response.Body.Close()
					result = string(bodyBytes)
					statuscode = response.StatusCode
				}
			}
		}
	} else {
		statuscode = 500
		result = "failure - executesql"
	}
	returnResult := fmt.Sprintf(`{"statuscode":%d,"body":{"result":"%s"}}`, statuscode, result)
	log.Println(returnResult)

	return returnResult, nil
}

func main() {
	lambda.Start(HandleRequest)
}
