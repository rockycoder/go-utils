package persistance

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"bitbucket.org/rockycoder/dataextractor/model"
	"cloud.google.com/go/pubsub"
)

type SellerPubSub struct {
	Client *pubsub.Client
	Topic  *pubsub.Topic
}

var err error

func InitializeSellerPubSub(projectID, topicName string) *SellerPubSub {

	ctx := context.Background()

	lPubSubCient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)

	}

	topic := lPubSubCient.Topic(topicName)

	if topic == nil {
		topic, err = lPubSubCient.CreateTopic(ctx, topicName)
		if err != nil {

			log.Fatalf("Failed to create topic: %v", err)

		}

	}

	return &SellerPubSub{Client: lPubSubCient, Topic: topic}

}

func (spb *SellerPubSub) PublishMessage(product model.ProductSchema) {

	ctx := context.Background()
	b, err := json.Marshal(product)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := &pubsub.Message{
		Data: b,
	}

	if _, err := spb.Topic.Publish(ctx, msg); err != nil {
		fmt.Println(err)
		return
	}

}

func (spb *SellerPubSub) Subscribe(subName string) {

	ctx := context.Background()

	subscription, _ := spb.Client.CreateSubscription(ctx, subName, spb.Topic, 0, nil)

	it, err := subscription.Pull(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg, err := it.Next()
		if err != nil {
			log.Fatalf("could not pull: %v", err)
		}
		var data model.ProductSchema
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Done(true)
			continue
		}

		fmt.Println("[ID %d] Processing.", data)

	}
}
