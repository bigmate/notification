package kafka

import (
	"emailservice/pkg/logger"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//TODO: make a repo
//getKakkaConsumer gets the consumer
func (s *service) getKafkaConsumer() (*kafka.Consumer, error) {
	//kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-4nym6.us-east-1.aws.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.username":     "EBATE3PEHXNWTAWK",
		"sasl.password":     "qkI+cwpe4QM3Xn+ZBFq/nm2l81eU/dnLbfTilneQmpHhT4dS3Q6EE80qjY82LSkg",
		"sasl.mechanism":    "PLAIN",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		//todo: add error log
		panic(err)
	}

	return consumer, err
}

func (s *service) KafkaGetActivationCode() (interface{}, error) {
	var activationCode interface{}
	//kafka consumer
	consumer, err := s.getKafkaConsumer()
	if err != nil {
		return "", err
	}
	//TODO: create new topic on confluent cloud for each microservice
	consumer.SubscribeTopics([]string{"Default"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			logger.Errorf("Kafka: Error %v %v\n", err, msg)
			return "", err
		}
		logger.Infof("Kafka: message on %s %s\n", msg.TopicPartition, string(msg.Value))

		var message map[string]interface{}

		json.Unmarshal(msg.Value, &message)

		activationCode = message["code"]

	}
	consumer.Close()

	return activationCode, nil
}

func (s *service) KafkaGetResetCode() (interface{}, error) {
	var pwResetCode interface{}
	//kafka consumer
	consumer, err := s.getKafkaConsumer()
	if err != nil {
		return "", err
	}
	//Todo: create new topic on confluent cloud for each microservice
	consumer.SubscribeTopics([]string{"Default"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			logger.Errorf("Kafka: Error %v %v\n", err, msg)
			return "", err
		}
		logger.Infof("Kafka: message on %s %s\n", msg.TopicPartition, string(msg.Value))

		var message map[string]interface{}

		json.Unmarshal(msg.Value, &message)

		pwResetCode = message["resetCode"]

	}
	consumer.Close()

	return pwResetCode, nil
}

func (s *service) KafkaGetUsersEmail() (interface{}, error) {
	var userEmail interface{}
	//kafka consumer
	consumer, err := s.getKafkaConsumer()
	if err != nil {
		return "", err
	}
	//Todo: create new topic on confluent cloud for each microservice
	consumer.SubscribeTopics([]string{"Default"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			logger.Errorf("Kafka: Error %v %v\n", err, msg)
			return "", err
		}
		logger.Infof("Kafka: message on %s %s\n", msg.TopicPartition, string(msg.Value))

		var message map[string]interface{}

		json.Unmarshal(msg.Value, &message)

		userEmail = message["email"]

	}
	consumer.Close()

	return userEmail, nil
}
