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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dunv/ulog"
	"github.com/dunv/webhookForwarder/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "a client for forwarding webhookCalls",
	RunE: func(cmd *cobra.Command, args []string) error {
		incomingSocket, err := cmd.Flags().GetString("incomingSocket")
		if err != nil {
			return err
		}

		outgoingURI, err := cmd.Flags().GetString("outgoingURI")
		if err != nil {
			return err
		}

		for {
		Connect:

			httpClient := &http.Client{}

			ulog.Infof("connecting to server (%s)", incomingSocket)
			conn, err := grpc.Dial(incomingSocket, grpc.WithInsecure())
			if err != nil {
				ulog.Errorf("could not connect (%s)", err)
				time.Sleep(time.Second)
				goto Connect
			}

			client := pb.NewWebhookForwarderClient(conn)
			stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{
				Origin: "*",
			})
			if err != nil {
				ulog.Errorf("could not subscribe (%s)", err)
				time.Sleep(time.Second)
				goto Connect
			}

			ulog.Infof("connected and subscribed")
			for {
				call, err := stream.Recv()
				if err != nil {
					ulog.Errorf("could not receive (%s)", err)
					time.Sleep(time.Second)
					goto Connect
				}

				dest := fmt.Sprintf("%s%s", outgoingURI, call.GetUri())
				r, err := http.NewRequest(call.GetMethod(), dest, bytes.NewReader(call.GetBody()))
				if err != nil {
					ulog.Errorf("could not forward (%s)", err)
					continue
				}

				// reset default go-header
				r.Header.Del("Accept-Encoding")

				for k, v := range call.GetHeaders() {
					for _, singleValue := range v.GetEntry() {
						r.Header.Add(k, singleValue)
					}
				}

				res, err := httpClient.Do(r)
				if err != nil {
					ulog.Errorf("could not forward (%s)", err)
					continue
				}

				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					ulog.Errorf("could not forward (%s)", err)
					continue
				}
				res.Body.Close()

				ulog.Infof("forwarded destination: %s status: %s, body: %s", dest, res.Status, strings.TrimSpace(string(resBody)))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("incomingSocket", "i", "0.0.0.0:50051", "grpc socket to connect to")
	clientCmd.Flags().StringP("outgoingURI", "o", "http://0.0.0.0:8081", "socket to to forwarder webhook calls to")
}
