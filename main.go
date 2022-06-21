package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/oklog/run"
	"golang.org/x/net/context"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	defaultDomain = "kove.net"
)

func main() {
	if err := startPlugin(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// func generateValue() string {
// 	cmd, err := exec.Command("/bin/sh", "kove-pool-capacity-utility/helloworld.sh").Output()
// 	if err != nil {
// 		fmt.Printf("error %s", err)
// 	}
// 	output := string(cmd)
// 	return output
// }

func startPlugin() error {

	domain := flag.String("domain", defaultDomain, "The domain to use when when declaring devices.")
	pluginPath := flag.String("plugin-directory", v1beta1.DevicePluginPath, "The directory in which to create plugin sockets.")

	fmt.Printf("Hello world!")

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var g run.Group
	{
		d := new(DeviceSpec)
		d.Name = path.Join(*domain, "/memory")
		d.Count = 333
		// d.Value = generateValue()
		ctx, cancel := context.WithCancel(context.Background())
		gp := NewGenericPlugin(d, *pluginPath, log.With(logger, "resource", d.Name))
		// Start the generic device plugin server.
		g.Add(func() error {
			logger.Log("msg", fmt.Sprintf("Starting the kove-k8s-device-plugin for %q.", d.Name))
			return gp.Run(ctx)
		}, func(error) {
			cancel()
		})
		fmt.Printf("count %d", d.Count)
	}
	return g.Run()
}
