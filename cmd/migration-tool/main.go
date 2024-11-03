/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-08-23 18:09:11
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-31 14:17:51
 * @FilePath: /DappsterOS/cmd/migration-tool/main.go
 * @Description:
 * @Website: https://www.dappsteros.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package main

import (
	"flag"
	"fmt"
	"os"

	interfaces "github.com/dappster-io/DappsterOS-Common"
	"github.com/dappster-io/DappsterOS-Common/utils/systemctl"
	"github.com/dappster-io/DappsterOS/common"
	"github.com/dappster-io/DappsterOS/pkg/config"
	"github.com/dappster-io/DappsterOS/pkg/sqlite"
	"github.com/dappster-io/DappsterOS/service"
	"gorm.io/gorm"
)

const (
	dappsterosServiceName = "dappsteros.service"
)

var (
	commit = "private build"
	date   = "private build"

	_logger  *Logger
	sqliteDB *gorm.DB

	configFlag = ""
	dbFlag     = ""
)

func init() {
	versionFlag := flag.Bool("v", false, "version")
	debugFlag := flag.Bool("d", true, "debug")
	forceFlag := flag.Bool("f", true, "force")

	flag.Parse()

	if *versionFlag {
		fmt.Println("v" + common.VERSION)
		os.Exit(0)
	}

	println("git commit:", commit)
	println("build date:", date)

	_logger = NewLogger()

	if os.Getuid() != 0 {
		_logger.Info("Root privileges are required to run this program.")
		os.Exit(1)
	}

	if *debugFlag {
		_logger.DebugMode = true
	}

	if !*forceFlag {
		serviceEnabled, err := systemctl.IsServiceEnabled(dappsterosServiceName)
		if err != nil {
			panic(err)
		}

		if serviceEnabled {
			_logger.Info("%s is already enabled. If migration is still needed, try with -f.", dappsterosServiceName)
			os.Exit(1)
		}
	}

	config.InitSetup(configFlag, "")

	if len(dbFlag) == 0 {
		dbFlag = config.AppInfo.DBPath + "/db"
	}

	sqliteDB = sqlite.GetDb(dbFlag)
	// gredis.GetRedisConn(config.RedisInfo),

	service.MyService = service.NewService(sqliteDB, "")
}

func main() {
	migrationTools := []interfaces.MigrationTool{
		// nothing to migrate from last version
	}

	var selectedMigrationTool interfaces.MigrationTool

	// look for the right migration tool matching current version
	for _, tool := range migrationTools {
		migrationNeeded, err := tool.IsMigrationNeeded()
		if err != nil {
			panic(err)
		}

		if migrationNeeded {
			selectedMigrationTool = tool
			break
		}
	}

	if selectedMigrationTool == nil {
		_logger.Error("selectedMigrationTool is null")
		return
	}

	if err := selectedMigrationTool.PreMigrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.Migrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.PostMigrate(); err != nil {
		_logger.Error("Migration succeeded, but post-migration failed: %s", err)
	}
}
