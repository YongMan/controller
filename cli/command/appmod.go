package command

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/ksarch-saas/cc/cli/context"
	"github.com/ksarch-saas/cc/meta"
)

var AppModCommand = cli.Command{
	Name:   "appmod",
	Usage:  "appmod",
	Action: appModAction,
	Flags: []cli.Flag{
		cli.StringFlag{"d,appname", "", "appname"},
		cli.StringFlag{"s,enableslaveread", "", "AutoEnableSlaveRead <true> or <false>"},
		cli.StringFlag{"m,enablemasterwrite", "", "AutoEnableMasterWrite <true> or <false>"},
		cli.StringFlag{"f,failover", "", "AutoFailover <true> or <false>"},
		cli.IntFlag{"i,interval", -1, "AutoFailoverInterval"},
		cli.StringFlag{"r,masterregion", "", "MasterRegion"},
		cli.StringFlag{"R,regions", "", "Regions"},
		cli.IntFlag{"k,migratekey", -1, "MigrateKeysEachTime"},
		cli.IntFlag{"p,migratekeystep", -1, "MigrateKeysStep"},
		cli.IntFlag{"t,migratetimeout", -1, "MigrateTimeout"},
		cli.StringFlag{"l,slavefailoverlimit", "", "Slave failover limit check"},
		cli.IntFlag{"u,fetchinterval", -1, "Fetch cluster nodes interval"},
		cli.IntFlag{"c,migrateconcurrency", -1, "number of migrate task running concurrently"},
		cli.IntFlag{"e,fixclustercircle", -1, "period of fix cluster , s(second)"},
		cli.StringFlag{"a,AotuFixCluster", "", "AotuFixCluster <true> or <false>"},
	},
	Description: `
    update app configuraton in zookeeper
    `,
}

func appModAction(c *cli.Context) {
	appname := c.String("d")
	if appname == "" {
		appname = context.GetAppName()
	}
	s := c.String("s")
	m := c.String("m")
	f := c.String("f")
	i := c.Int("i")
	r := c.String("r")
	R := c.String("R")
	k := c.Int("k")
	p := c.Int("p")
	t := c.Int("t")
	l := c.String("l")
	u := c.Int("u")
	mc := c.Int("c")
	e := c.Int("e")
	a := c.String("a")

	appConfig := meta.AppConfig{}
	config, version, err := context.GetApp(appname)
	fmt.Println(string(config))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(config, &appConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	//update config if needed
	if s != "" {
		if s == "true" {
			appConfig.AutoEnableSlaveRead = true
		} else if s == "false" {
			appConfig.AutoEnableSlaveRead = false
		}
	}
	if m != "" {
		if m == "true" {
			appConfig.AutoEnableMasterWrite = true
		} else if m == "false" {
			appConfig.AutoEnableMasterWrite = false
		}
	}
	if f != "" {
		if f == "true" {
			appConfig.AutoFailover = true
		} else if f == "false" {
			appConfig.AutoFailover = false
		}
	}
	if l != "" {
		if l == "true" {
			appConfig.SlaveFailoverLimit = true
		} else if l == "false" {
			appConfig.SlaveFailoverLimit = false
		}
	}
	if i != -1 {
		appConfig.AutoFailoverInterval = time.Duration(i)
	}
	if r != "" {
		appConfig.MasterRegion = r
	}
	if R != "" {
		appConfig.Regions = strings.Split(R, ",")
	}
	if a != "" {
		if a == "true" {
			appConfig.AotuFixCluster = true
		} else if a == "false" {
			appConfig.AotuFixCluster = false
		}
	}
	if k != -1 {
		appConfig.MigrateKeysEachTime = k
	}
	if p != -1 {
		appConfig.MigrateKeysStep = p
	}
	if t != -1 {
		appConfig.MigrateTimeout = t
	}
	if u != -1 && u > 0 {
		appConfig.FetchClusterNodesInterval = time.Duration(u)
	}
	if mc != -1 {
		appConfig.MigrateConcurrency = mc
	}
	if e != -1 {
		appConfig.FixClusterCircle = e
	}

	out, err := json.Marshal(appConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = context.ModApp(appname, out, version)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Mod %s success\n%s\n", appname, string(out))
}
