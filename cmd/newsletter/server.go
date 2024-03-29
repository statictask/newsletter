// Copyright © 2022 Luan Guimarães Lacerda <luang@riseup.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/config"
	"github.com/statictask/newsletter/pkg/scheduler"
	"github.com/statictask/newsletter/pkg/server"
	"go.uber.org/zap"
)

// serverCmd represents the start command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "newsletter server API",
	Long:  `newsletter server API`,
	Run:   startServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("bind", "", "server bind address")
}

func startServer(cmd *cobra.Command, args []string) {
	log.L.Info("initializing system")

	initDB(cmd, args)
	initSchedulers(cmd, args)
	initServer(cmd, args)

	log.L.Info("finished")
}

func initDB(cmd *cobra.Command, args []string) {
	database.Init()

	if err := database.Ping(); err != nil {
		log.L.Fatal("failed connecting to postgres", zap.Error(err))
	}

	log.L.Info("successfully connected to postgres!")
}

func initSchedulers(cmd *cobra.Command, args []string) {
	ps := scheduler.NewPipelineScheduler()
	ps.Start()

	ts := scheduler.NewTaskScheduler()
	ts.Start()

	ss := scheduler.NewScrapperJobScheduler()
	ss.Start()

	sm := scheduler.NewPublisherJobScheduler()
	sm.Start()
}

func initServer(cmd *cobra.Command, args []string) {
	bind, err := cmd.Flags().GetString("bind")
	if err != nil {
		log.L.Fatal("option --bind is missing", zap.Error(err))
	}

	if bind == "" {
		bind = config.C.BindAddress
	}

	s := server.New()

	if err := s.Listen(bind); err != nil {
		log.L.Fatal("unable to start server", zap.Error(err))
	}
}
