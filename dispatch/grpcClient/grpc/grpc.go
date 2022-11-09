package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/team-triage/triage/dispatch/grpcClient/pb" // import protobuf module

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func makeConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    time.Second * 3, // how long we wait to hear back from the server before closing connection
		Timeout: time.Second * 1, // frequency of pings
	}))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn, err
}

func ConnectToServer(address string) pb.MessageHandlerClient {
	conn, err := makeConnection(address)
	if err != nil {
		fmt.Println("GRPC: We got an error", err)
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewMessageHandlerClient(conn) // init client

	return client
	// defer conn.Close()
}

func SendMessage(client pb.MessageHandlerClient, msgValue string) (int32, error) { // will update parameter from string to proper struct
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5) // nice to have: adjust this and make it configurable during deployment

	defer cancel()

	fmt.Println("GRPC CLIENT: Sending message!", msgValue)

	resp, err := client.SendMessage(ctx, &pb.Message{Body: msgValue})

	fmt.Println(resp)

	if err != nil {
		return int32(0), err
		// return zero-valued int32, error
		// log.Fatalf("could not get message: %v", err)
		// push message to messages channel, then break out of wrapping goRoutine
	}

	return resp.GetStatus(), nil // "ack" or "nack"
}
