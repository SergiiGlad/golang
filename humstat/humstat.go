package humstat

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"go-team-room/humaws" //Need it for DynamoDb session obj
)

// SendStat - chanel there to send your statistic metrics
var SendStat chan map[string]int

// OutputStat - chanel from where to get your statistic metrics
var OutputStat chan map[string]int

func init() {
	SendStat = make(chan map[string]int, 10000)
	OutputStat = make(chan map[string]int, 0)

	go sumStat(SendStat, OutputStat)

}

//sumStat - funk to count total sum of all metrics
func sumStat(inputChan chan map[string]int, outputChan chan map[string]int) {
	var someMetric, someRequest map[string]int
	TotalStat := make(map[string]int)
	AppStartTime := time.Now().Unix()
	MinuteStat := make(map[string]int)
	//OneMinuteTimeCounter := time.After(time.Minute)
	OneMinuteTimeCounter := time.Tick(time.Minute)
	/* начать рефарториг с этого цыкла! */
	/* ОЧЕНЬ не эффективно              */
	for {
		select {
		case someMetric = <-inputChan:
			//Some data come to us
			//fmt.Println(TotalStat)
			for key, value := range someMetric {
				if _, ok := MinuteStat[key]; ok {
					MinuteStat[key] = MinuteStat[key] + value
				} else {
					MinuteStat[key] = value
				}
			}
			//fmt.Println(MinuteStat)

		case <-OneMinuteTimeCounter:
			//Save one minute stat to DynamoDB
			av, err := dynamodbattribute.MarshalMap(MinuteStat)
			fmt.Println(av) ///////////////////
			if err != nil {
				fmt.Println("Got error marshalling map:")
				fmt.Println(err.Error())
				//os.Exit(1)
			}
			av["timestamp"], _ = dynamodbattribute.Marshal(time.Now().Format("2006-01-02T15:04"))
			//av["id_timestamp"], _ = dynamodbattribute.Marshal(fmt.Sprint(time.Now().UnixNano()))
			av["id_timestamp"], _ = dynamodbattribute.Marshal(time.Now().Format("2006-01-02"))

			// Create item in table messages
			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String("humstat"),
			}
			/*
				<script src="https://sdk.amazonaws.com/js/aws-sdk-2.195.0.min.js"></script>
			*/
			_, err = humaws.Dyna.Db.PutItem(input)

			if err != nil {
				fmt.Println("Got error calling PutItem:")
				fmt.Println(err.Error())
				//os.Exit(1)

			}
			for j := range MinuteStat {
				if _, ok := TotalStat[j]; ok {
					TotalStat[j] = TotalStat[j] + MinuteStat[j]
				} else {
					TotalStat[j] = MinuteStat[j]
				}
				MinuteStat[j] = 0
			}
			//OneMinuteTimeCounter = time.After(time.Minute)
		case someRequest = <-outputChan:
			//some one trying to get stat
			for key, value := range someRequest {
				if key == "stat" && value == 0 {
					TotalStat["uptime_seconds"] = int(time.Now().Unix() - AppStartTime)
					outputChan <- TotalStat
					delete(TotalStat, "uptime_seconds")
				}
				if key == "stat" && value == 1 {
					outputChan <- MinuteStat
				}
			}
		default:
			//тут можно ничего не писать, чтобы данные молча отбрасывались
			//fmt.Println("потрачено")
		}
		time.Sleep(time.Millisecond * 5)
	}
}
