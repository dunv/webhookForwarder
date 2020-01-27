package cmd

import (
	"log"
	"net"
	"time"

	"github.com/dunv/connectionTools"
	"github.com/dunv/ulog"
	"github.com/dunv/webhookForwarder/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer(outgoingSocket string, hub *connectionTools.NotificationHub) {
	lis, err := net.Listen("tcp", outgoingSocket)
	if err != nil {
		log.Printf("could not listen %v \n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterWebhookForwarderServer(s, &GrpcServer{
		hub: hub,
	})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		ulog.Fatalf("could not serve %v", err)
	}

}

type GrpcServer struct {
	hub *connectionTools.NotificationHub
}

func (s GrpcServer) Subscribe(req *pb.SubscribeRequest, srv pb.WebhookForwarder_SubscribeServer) error {
	callChannel := make(chan interface{})
	connGUID := s.hub.Register("all", callChannel)
	ulog.Infof("client connected (%s)", connGUID)

	for {
		select {
		case call := <-callChannel:

			err := srv.Send(call.(*pb.WebhookCall))
			if err != nil {
				s.hub.Unregister(connGUID, err)
				ulog.Infof("client disconnected (%s)", connGUID)
				return err
			}
		case <-time.After(5 * time.Second):
			err := srv.Context().Err()
			if err != nil {
				s.hub.Unregister(connGUID, err)
				ulog.Infof("client disconnected (%s)", connGUID)
				return err
			}
		}
	}
}
