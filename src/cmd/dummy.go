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
	"net/http"
	"net/http/httputil"

	"github.com/dunv/uhttp"
	"github.com/spf13/cobra"
)

// dummyCmd represents the dummy command
var dummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "dummy",
	RunE: func(cmd *cobra.Command, args []string) error {
		incomingSocket, err := cmd.Flags().GetString("incomingSocket")
		if err != nil {
			return err
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := httputil.DumpRequest(r, true)
			if err != nil {
				uhttp.RenderError(w, r, err)
			}
			fmt.Println(string(body))
			fmt.Println()
			uhttp.Render(w, r, map[string]string{"all": "good"})
		})

		return http.ListenAndServe(incomingSocket, nil)
	},
}

func init() {
	rootCmd.AddCommand(dummyCmd)
	dummyCmd.Flags().StringP("incomingSocket", "i", "0.0.0.0:8081", "exposed socket for logging")
}
