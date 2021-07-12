/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ryicoh/apery-graphql/pkg/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "apery-graphql",
	Short: "",
	Long:  "",
	RunE:  run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var (
	port   int
	binary string
)

func init() {
	rootCmd.Flags().IntVar(&port, "port", 9012, "server port")
	rootCmd.Flags().StringVar(&binary, "binary", "apery", "binary path")

}

func run(cmd *cobra.Command, args []string) (err error) {
	if p := os.Getenv("PORT"); p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	}

	srv := server.NewServer(port, binary)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		fmt.Printf("apery-graphql listen :%d\n", port)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-errCh:
		return err
	}

	return nil
}
