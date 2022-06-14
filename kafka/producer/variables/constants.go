package variables

import (
	_ "github.com/joho/godotenv/autoload"
)

var KafkaBootstrapServers = "localhost:29092"
var KafkaGroupId = "go-kafka-simple"
var KafkaTopic = "go-kafka-simple"
var SentryLink = "https://f7413278bc0f4d95a1d4e1a164ec8b03@o204768.ingest.sentry.io/5880501"
var NumCore = 1

//var KafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
//var KafkaGroupId = os.Getenv("KAFKA_GROUP_ID")
//var KafkaTopic = os.Getenv("KAFKA_TOPIC")
//var SentryLink = os.Getenv("SENTRY_LINK")
//var NumCore, _ = strconv.Atoi(os.Getenv("NUMBER_CORE"))
