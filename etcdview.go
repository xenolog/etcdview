// entry point to the Janus ()
package main

import (
	"errors"
	"fmt"
	etcdc "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/xenolog/go-tiny-logger.v1"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	Version = "0.0.1"
)

var (
	Log *logger.Logger
	App *cli.App
	err error
)

func init() {
	// Setup logger
	Log = logger.New()

	// Configure CLI flags and commands
	App = cli.NewApp()
	App.Name = "ETCD subtree listener"
	App.Version = Version
	App.EnableBashCompletion = true
	App.Usage = "Specify entry point of tree and got subtree for simple displaying"
	App.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug mode. Show more output",
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: "http://127.0.0.1:4001",
			Usage: "Specify URL for connect to ETCD",
		},
	}
	App.Commands = []cli.Command{{
		Name:   "tree",
		Usage:  "Display tree (or subtree) for etcd",
		Action: displayTree,
		// }, {
		// 	Name:   "test",
		// 	Usage:  "Just run smaLog test.",
		// 	Action: runTest,
	},
	}
	App.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			Log.SetMinimalFacility(logger.LOG_D)
		} else {
			Log.SetMinimalFacility(logger.LOG_I)
		}
		Log.Debug("EtcdTree started.")
		return nil
	}
	App.CommandNotFound = func(c *cli.Context, cmd string) {
		Log.Printf("Wrong command '%s'", cmd)
		os.Exit(1)
	}
}

func main() {
	App.Run(os.Args)
}

func getHeader(startKey string) (string, int) {
	var rv string
	level := 0
	s := strings.Trim(startKey, "/")
	ss := strings.Split(s, "/")
	for _, v := range ss[:len(ss)-1] {
		rv += fmt.Sprintf("%s%s:\n", strings.Repeat("  ", level), v)
		level += 1
	}
	return rv, level
}

func displaySubtree(node *etcdc.Node, level int) {
	keys := strings.Split(node.Key, "/")
	currKey := keys[len(keys)-1]
	currPrefix := strings.Repeat("  ", level)
	if node.Dir {
		// process a subtree
		Log.Debug("'%q' key is a subtree\n", node.Key)
		fmt.Printf("%s%s:\n", currPrefix, currKey)
		for _, v := range node.Nodes {
			displaySubtree(v, level+1)
		}
	} else {
		// just a key/value
		fmt.Printf("%s%s: %s\n", currPrefix, currKey, node.Value)
	}
}

func displayTree(c *cli.Context) error {
	root := c.Args()[0]
	Log.Debug("Display tree from '%s'", root)
	// create transport for connect to etcd
	etcdTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	// connect to etcd
	endpoints := []string{c.GlobalString("url")}
	Log.Debug("ETCD endpoints: '%s'", strings.Join(endpoints, ":"))
	etcdClient, err := etcdc.New(etcdc.Config{
		Endpoints:               endpoints,
		Transport:               etcdTransport,
		HeaderTimeoutPerRequest: time.Second,
	})
	if err != nil {
		errmsg := fmt.Sprintf("Can't connect to etcd by endpoints: '%s'", strings.Join(endpoints, ","))
		Log.Error(errmsg)
		return errors.New(errmsg)
	}
	etcdAPI := etcdc.NewKeysAPI(etcdClient)
	// get
	resp, err := etcdAPI.Get(context.Background(), root, &etcdc.GetOptions{
		Recursive: true,
		Sort:      true,
	})
	if err != nil {
		Log.Error("Error: %s", err)
	} else {
		// print common key info
		Log.Debug("Get is done. Metadata is %q\n", resp)
		header, level := getHeader(root)
		fmt.Print(header)
		displaySubtree(resp.Node, level)
	}
	return nil
}
