package app

import (
	"github.com/ezaurum/cthulthu/context"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/node"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
)

func initialize() context.Context {
	var ctx = context.Ctx()
	if dbg := os.Getenv("SS_DEBUG"); len(dbg) > 0 && strings.ToLower(dbg) != "false" {
		ctx.SetDebug(true)
	}
	// get node number
	nodeNumber := node.ByIP()
	log.Printf("node number: %d", nodeNumber)
	ctx.SetNodeNumber(nodeNumber)

	sn := snowflake.New(nodeNumber)
	var idGenerators generators.IDGenerators
	var migrates []interface{}
	var merr error
	if migrates, merr = ctx.ResourceInterfaces(); nil != merr {
		log.Fatal("migrates error", merr)
	}
	idGenerators = generators.New(func(typeString string) generators.IDGenerator {
		return sn
	}, migrates...)
	ctx.SetIDGenerators(idGenerators)

	if repo, err := database.New(idGenerators, ctx.Debug()); nil != err {
		log.Fatal("database init error", err.Error())
	} else {
		log.Info("masterDB - ", repo.Writer())
		log.Info("slaveDB - ", repo.Reader())
		repo.Writer().AutoMigrate(migrates...)
		ctx.SetRepository(repo)
	}
	return ctx
}
