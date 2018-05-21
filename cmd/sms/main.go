package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kgolding/hilink"
)

var (
	flagEndpoint = flag.String("endpoint", "http://192.168.1.1/", "api endpoint")
	flagDebug    = flag.Bool("v", false, "enable verbose")
	flagTo       = flag.String("to", "", "to")
	flagMsg      = flag.String("msg", "", "message")
	flagList     = flag.Bool("list", false, "list sms messages in inbox")
	flagCount    = flag.Uint("c", 50, "message count for -list")
	flagNoLogin  = flag.Bool("nologin", false, "do not try to login")
	flagLogin    = flag.String("login", "admin", "admin")
	flagPass     = flag.String("pass", "admin", "admin")
)

func main() {
	var err error

	flag.Parse()

	// options
	opts := []hilink.Option{
		hilink.URL(*flagEndpoint),
	}
	if *flagDebug {
		opts = append(opts, hilink.Log(log.Printf, log.Printf))
	}

	if !*flagNoLogin {
		opts = append(opts, hilink.Login(*flagLogin, *flagPass))
	}

	// create client
	client, err := hilink.NewClient(opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// handle list
	if *flagList {
		doList(client, hilink.SmsBoxTypeInbox, *flagCount)
		return
	}

	// check flags
	if *flagMsg == "" {
		fmt.Fprintf(os.Stderr, "error: must specify msg\n")
		os.Exit(1)
	}
	if *flagTo == "" {
		fmt.Fprintf(os.Stderr, "error: must specify to\n")
		os.Exit(1)
	}

	// send sms
	b, err := client.SmsSend(*flagMsg, *flagTo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if !b {
		fmt.Fprintf(os.Stderr, "could not send message\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "message sent\n")
}

// doList lists the sms in the inbox in json format.
func doList(client *hilink.Client, bt hilink.SmsBoxType, count uint) {
	// get sms counts
	l, err := client.SmsList(uint(bt), 1, count, false, false, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// convert to json
	buf, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(buf))
}
