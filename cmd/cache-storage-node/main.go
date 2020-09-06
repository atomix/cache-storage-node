// Copyright 2020-present Open Networking Foundation.
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
	"bytes"
	"fmt"
	"github.com/atomix/api/proto/atomix/database"
	"github.com/atomix/go-framework/pkg/atomix/counter"
	"github.com/atomix/go-framework/pkg/atomix/election"
	"github.com/atomix/go-framework/pkg/atomix/indexedmap"
	"github.com/atomix/go-framework/pkg/atomix/leader"
	"github.com/atomix/go-framework/pkg/atomix/list"
	logprimitive "github.com/atomix/go-framework/pkg/atomix/log"
	"github.com/atomix/go-framework/pkg/atomix/map"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	"github.com/atomix/go-framework/pkg/atomix/set"
	"github.com/atomix/go-framework/pkg/atomix/value"
	"io/ioutil"
	"net"
	"os"
	"os/signal"

	"github.com/atomix/go-local/pkg/atomix/local"
	"github.com/gogo/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stdout)

	clusterConfig := parseClusterConfig()

	// Start the node. The node will be started in its own goroutine.
	member := clusterConfig.Replicas[0]
	log.Info(member.Host, ":", member.APIPort)
	address := fmt.Sprintf(":%d", member.APIPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a local node
	partitions := make([]primitive.PartitionID, len(clusterConfig.Partitions))
	for i, partition := range clusterConfig.Partitions {
		partitions[i] = primitive.PartitionID(partition.Partition)
	}
	node := local.NewNode(lis, partitions)

	// Register primitives on the Atomix node
	counter.RegisterPrimitive(node)
	election.RegisterPrimitive(node)
	indexedmap.RegisterPrimitive(node)
	logprimitive.RegisterPrimitive(node)
	leader.RegisterPrimitive(node)
	list.RegisterPrimitive(node)
	_map.RegisterPrimitive(node)
	set.RegisterPrimitive(node)
	value.RegisterPrimitive(node)

	// Start the node
	if err := node.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for an interrupt signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop the node after an interrupt
	if err := node.Stop(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseClusterConfig() *database.DatabaseConfig {
	clusterConfigFile := os.Args[2]
	clusterConfig := &database.DatabaseConfig{}
	clusterConfigBytes, err := ioutil.ReadFile(clusterConfigFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := jsonpb.Unmarshal(bytes.NewReader(clusterConfigBytes), clusterConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return clusterConfig
}
