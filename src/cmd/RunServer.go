/*
Copyright © 2020 Daniel Unverricht <daniel@unvericht.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
