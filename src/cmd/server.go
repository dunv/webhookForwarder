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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/dunv/connectionTools"
	"github.com/dunv/uhttp"
	"github.com/dunv/ulog"
	"github.com/dunv/webhookForwarder/pb"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts a server for receiving webhook calls",
	RunE: func(cmd *cobra.Command, args []string) error {
		incomingSocket, err := cmd.Flags().GetString("incomingSocket")
		if err != nil {
			return err
		}

		incomingPath, err := cmd.Flags().GetString("incomingPath")
		if err != nil {
			return err
		}

		outgoingSocket, err := cmd.Flags().GetString("outgoingSocket")
		if err != nil {
			return err
		}

		printDump, err := cmd.Flags().GetBool("printHttpDump")
		if err != nil {
			return err
		}

		hub := connectionTools.NewNotificationHub(5 * time.Second)

		go RunServer(outgoingSocket, hub)

		http.HandleFunc(incomingPath, func(w http.ResponseWriter, r *http.Request) {
			body, err := httputil.DumpRequest(r, true)
			if err != nil {
				uhttp.RenderError(w, r, err)
			}

			if printDump {
				ulog.Info("received webhook call")
				fmt.Println()
				fmt.Println(string(body))
				fmt.Println()
			}

			headers := map[string]*pb.ListOfString{}
			for k, v := range r.Header {
				headers[k] = &pb.ListOfString{Entry: v}
			}

			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				uhttp.RenderError(w, r, err)
				return
			}

			call := &pb.WebhookCall{
				Origin:  r.RemoteAddr,
				Method:  r.Method,
				Uri:     r.RequestURI,
				Headers: headers,
				Body:    body,
			}

			result, errMap := hub.Notify("all", call)
			ulog.Infof("forwarded webhook call (%d recipients)", result)
			uhttp.Render(w, r, struct {
				SuccessfulNotifications int              `json:"successfulNotifications"`
				ErrMap                  map[string]error `json:"notifyErrors,omitempty"`
			}{
				SuccessfulNotifications: result,
				ErrMap:                  errMap,
			})

		})

		ulog.Infof("Server running on %s", incomingSocket)
		return http.ListenAndServe(incomingSocket, nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("incomingSocket", "i", "0.0.0.0:8080", "socket to listen on for webhook calls")
	serverCmd.Flags().StringP("incomingPath", "p", "/", "path to listen to for websocket calls")
	serverCmd.Flags().StringP("outgoingSocket", "o", "0.0.0.0:50051", "socket to listen on for client-connections")
	serverCmd.Flags().BoolP("printHttpDump", "d", false, "dump every webhook call to stdout")
}
