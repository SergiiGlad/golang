package humstat

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"go-team-room/humaws" //Need it fo DynamoDb session obj
)

// SendStat - chanel there to send your statistic metrics
var SendStat chan map[string]int

// OutputStat - chanel from where to get your statistic metrics
var OutputStat chan map[string]int

func init() {
	SendStat = make(chan map[string]int, 10000)
	OutputStat = make(chan map[string]int, 10)

	go sumStat(SendStat, OutputStat)

}

//sumStat - funk to count total sum of all metrics
func sumStat(inputChan chan map[string]int, outputChan chan map[string]int) {
	var someMetric map[string]int
	TotalStat := make(map[string]int)
	MinuteStat := make(map[string]int)

	/* начать рефарториг с этого цыкла! */
	/* ОЧЕНЬ не єффективно              */
	for {
		select {
		case someMetric = <-inputChan:
			//Some data come to us
			for key, value := range someMetric {
				if _, ok := TotalStat[key]; ok {
					TotalStat[key] = TotalStat[key] + value
				} else {
					TotalStat[key] = value
				}
			}
			//fmt.Println(TotalStat)
			for key, value := range someMetric {
				if _, ok := MinuteStat[key]; ok {
					MinuteStat[key] = MinuteStat[key] + value
				} else {
					MinuteStat[key] = value
				}
			}
			//fmt.Println(MinuteStat)

		case <-time.After(time.Minute):
			//Save one minute stat to DynamoDB
			av, err := dynamodbattribute.MarshalMap(MinuteStat)

			if err != nil {
				fmt.Println("Got error marshalling map:")
				fmt.Println(err.Error())
				//os.Exit(1)
			}

			// Create item in table messages
			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String("humstat"),
			}

			_, err = humaws.Dyna.Db.PutItem(input)

			if err != nil {
				fmt.Println("Got error calling PutItem:")
				fmt.Println(err.Error())
				//os.Exit(1)
			}
		default:
			//тут можно ничего не писать, чтобы данные молча отбрасывались
			//fmt.Println("потрачено")
		}
		time.Sleep(time.Millisecond * 5)
	}
}
