package main
/* 
  Copyright 2020 Google LLC. This software is provided as-is, without warranty or representation for any use or purpose. 
  Your use of it is subject to your agreement with Google.  
*/

import (
	"context"
	"flag"
	"fmt"
	"log/syslog"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

var pubsubTopic *string
var projectId *string
var syslogTag *string
var syslogServer *string
var syslogPort *int

func main() {
	pubsubTopic = flag.String("topic", "", "Pub/Sub topic")
	projectId = flag.String("project", os.Getenv("GOOGLE_PROJECT"), "Project ID")
	syslogTag := flag.String("tag", "pubsub2syslog", "Syslog tag")
	syslogServer := flag.String("server", "", "Syslog server hostname (host:port)")
	syslogProtocol := flag.String("protocol", "tcp", "Syslog protocol")
	syslogPriority := flag.Int("priority", int(syslog.LOG_NOTICE|syslog.LOG_DAEMON), "Syslog priority")

	flag.Parse()
	if *pubsubTopic == "" || *syslogTag == "" || *syslogServer == "" {
		fmt.Fprintf(os.Stderr, "Pub/Sub to syslog")
		flag.PrintDefaults()
		os.Exit(1)
	}

	writer, err := syslog.Dial(*syslogProtocol, *syslogServer, syslog.Priority(*syslogPriority), *syslogTag)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		panic(err)
	}

	topic := pubsubClient.Topic(*pubsubTopic)
	subscriptionName := fmt.Sprintf("%s-subscription", *syslogTag)
	sub, err := pubsubClient.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		sub = pubsubClient.Subscription(subscriptionName)
	}

	var mu sync.Mutex
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Fprintf(writer, "%q", string(msg.Data))
		fmt.Printf("%q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
	})
	if err != nil {
		panic(err)
	}
}
