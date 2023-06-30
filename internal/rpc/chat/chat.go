// Copyright © 2023 OpenIM open source community. All rights reserved.
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

package chat

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	chat2 "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	chatClient "github.com/OpenIMSDK/chat/pkg/rpclient/chat"
	"github.com/OpenIMSDK/chat/pkg/rpclient/openim"
	"github.com/OpenIMSDK/chat/pkg/sms"
	"google.golang.org/grpc"
)

func Start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := dbconn.NewGormDB()
	if err != nil {
		return err
	}
	tables := []any{
		chat2.Account{},
		chat2.Register{},
		chat2.Attribute{},
		chat2.VerifyCode{},
		chat2.UserLoginRecord{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}
	s, err := sms.New()
	if err != nil {
		return err
	}
	chat.RegisterChatServer(server, &chatSvr{
		Database: database.NewChatDatabase(db),
		Admin:    chatClient.NewAdminClient(discov),
		OpenIM:   openim.NewOpenIMClient(discov),
		SMS:      s,
	})
	return nil
}

type chatSvr struct {
	Database database.ChatDatabaseInterface
	Admin    *chatClient.AdminClient
	OpenIM   *openim.OpenIMClient
	SMS      sms.SMS
}
