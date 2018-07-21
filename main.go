package main

import (
	"io/ioutil"
	"fmt"
	"github.com/shamx9ir/gohue"
	"os"
	"log"
	"time"
	"strings"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

func main() {

	bridgesOnNetwork, _ := hue.FindBridges()
	bridge := bridgesOnNetwork[0]

	// read user from txt
	userFile, err := ioutil.ReadFile("user.txt")
	var userAuthenticated = false

	if err == nil {
		var userName = string(userFile)
		if userName != "" {
			err = bridge.Login(userName)
			if err != nil {
				fmt.Println("Failed to login as the user in user.txt.")
			} else {
				userAuthenticated = true
			}
		}
	} else 	{
		fmt.Println("Failed to open the file user.txt.")
	}

	if !userAuthenticated {
		fmt.Println("Please press button on the Hue bridge to login as new user. Then press the Enter Key to continue.")
		var input string
		fmt.Scanln(&input)

		username, _ := bridge.CreateUser("hueSensorCollector")
		bridge.Login(username)

		file, err := os.Create("user.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()

		fmt.Fprintf(file, username)
	}


	// get all sensors
	sensors, _ := bridge.GetAllSensors()

	// create name map
	sensorNameMap := make(map[string]string)
	for _, sensor := range sensors {
		// get sensor name from ZZPresence
		if sensor.Type == "ZLLPresence" {
			sensorId := strings.Split(sensor.UniqueID, "-")[0]
			sensorNameMap[sensorId] = sensor.Name
		}
	}

	// setup google client
	var projectID = ""
	projectFile, err := ioutil.ReadFile("gcpproject.txt")
	if err == nil {
		projectID = string(projectFile)
		if projectID == "" {
			fmt.Println("gcpproject.txt is empty.")
		}
	} else 	{
		fmt.Println("Failed to open the file gcpproject.txt.")
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID, option.WithCredentialsFile("gcp.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	for  { 	// do read every 5 mins
		for _, sensor := range sensors {
			if sensor.Type == "ZLLTemperature" {
				fmt.Println(sensor.Name, ":", float32(sensor.State.Temperature) / 100, "C")
				sensorId := strings.Split(sensor.UniqueID, "-")[0]
				LogTempReading(ctx, client, sensorNameMap[sensorId], sensorId, sensor.State.Temperature)
			}
		}

		time.Sleep(time.Minute * 5)
	}
}

func LogTempReading(ctx context.Context, client *datastore.Client, name string, uniqueId string, reading uint16) {

	logTime := time.Now()
	// Sets the kind for the new entity.
	kind := "TempLog"
	// Sets the name/ID for the new entity.
	keyName := uniqueId + "-" + fmt.Sprint(logTime)
	// Creates a Key instance.
	logKey := datastore.NameKey(kind, keyName, nil)

	tempLog := TempLog{
		SensorName: name,
		SensorId: uniqueId,
		TempValue: int16(reading),
		ValueDate: logTime,
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, logKey, &tempLog); err != nil {
		log.Fatalf("Failed to save TempLog: %v", err)
	}


}

type TempLog struct {
	SensorName string
	SensorId string
	TempValue int16
	ValueDate time.Time
}