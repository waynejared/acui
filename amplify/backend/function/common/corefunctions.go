package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret() (string, string) {
	secretName := "roxwen-postgres"
	region := "us-west-2"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	// Your code goes here.

	var dbCreds DbCreds
	json.Unmarshal([]byte(secretString), &dbCreds)
	return dbCreds.GetConnString(), dbCreds.GetEngine()
}

// This awesome f**ing function takes a sql string, an array of arguments and a struct
// It runs the sql string with the args and returns the object
// This expects a sql string to execute, 0-10 arguments of any type and a struct to
// unmarshal the result into
func ExecuteSQL(sqlString string, queryArgs []any, returnObj *ManagerInfoResponse) {
	connInfo, engine := GetSecret()

	//	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=enable", host, port, user, password, dbname)
	log.Println("Running ExecuteSQL")
	Conn, err := sql.Open(engine, connInfo)
	if err != nil {
		log.Println(err.Error())
	} else {
		var row *sql.Row
		defer Conn.Close()
		//sqlStmt := `select * from eac.addmanager($1);`
		switch len(queryArgs) {
		case 1:
			row = Conn.QueryRow(sqlString, queryArgs[0])
		case 2:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1])
		case 3:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2])
		case 4:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3])
		case 5:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4])
		case 6:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4], queryArgs[5])
		case 7:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4], queryArgs[5], queryArgs[6])
		case 8:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4], queryArgs[5], queryArgs[6], queryArgs[7])
		case 9:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4], queryArgs[5], queryArgs[6], queryArgs[7], queryArgs[8])
		case 10:
			row = Conn.QueryRow(sqlString, queryArgs[0], queryArgs[1], queryArgs[2], queryArgs[3], queryArgs[4], queryArgs[5], queryArgs[6], queryArgs[7], queryArgs[8], queryArgs[9])
		default:
			row = Conn.QueryRow(sqlString)
		}
		switch err := row.Scan(&returnObj.Body.Payload); err {
		case sql.ErrNoRows:
			fmt.Println("ExecuteSQL failed with: " + err.Error())
		case nil:
			returnObj.StatusCode = 200

		}
	}
}
