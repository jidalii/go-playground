package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ClusterArray []int

func (arr *ClusterArray) Set(val string) error {
	if len(*arr) > 0 {
		return errors.New("cluster flag already set")
	}
	for _, val := range strings.Split(val, ",") {
		cluster, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*arr = append(*arr, cluster)
	}
	return nil
}

func (arrs *ClusterArray) String() string {
	str := "["
	for _, s := range *arrs {
		str += strconv.Itoa(s)
	}
	str += "]"
	return str
}

var (
	host       string
	port       int
	isProdMode bool
	timeout    time.Duration
	clusters   ClusterArray
)

func main() {

	flag.StringVar(&host, "host", "localhost", "db host")
	flag.IntVar(&port, "p", 3306, "db port (shorthand)")
	flag.IntVar(&port, "port", 3306, "db port")
	flag.BoolVar(&isProdMode, "isProd", false, "isProdMode")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "db timeout")
	flag.Var(&clusters, "clusters", "db clusters")

	flag.Parse()

	fmt.Println("host:", host)
	fmt.Println("port:", port)
	fmt.Println("isProdMode:", isProdMode)
	fmt.Println("timeout:", timeout)
	fmt.Println("clusters:", clusters)
}
