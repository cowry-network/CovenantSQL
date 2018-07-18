/*
 * Copyright 2018 The ThunderDB Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package route

import (
	"fmt"

	"gitlab.com/thunderdb/ThunderDB/conf"
	"gitlab.com/thunderdb/ThunderDB/consistent"
	"gitlab.com/thunderdb/ThunderDB/crypto/kms"
	"gitlab.com/thunderdb/ThunderDB/proto"
	"gitlab.com/thunderdb/ThunderDB/utils/log"
)

// DHTService is server side RPC implementation
type DHTService struct {
	Consistent *consistent.Consistent
}

// NewDHTServiceWithRing will return a new DHTService and set an existing hash ring
func NewDHTServiceWithRing(c *consistent.Consistent) (s *DHTService, err error) {
	s = &DHTService{
		Consistent: c,
	}
	return
}

// NewDHTService will return a new DHTService
func NewDHTService(DHTStorePath string, persistImpl consistent.Persistence, initBP bool) (s *DHTService, err error) {
	c, err := consistent.InitConsistent(DHTStorePath, persistImpl, initBP)
	if err != nil {
		log.Errorf("init DHT service failed: %s", err)
		return
	}
	return NewDHTServiceWithRing(c)
}

// FindNode RPC returns node with requested node id from DHT
func (DHT *DHTService) FindNode(req *proto.FindNodeReq, resp *proto.FindNodeResp) (err error) {
	if !IsPermitted(&req.Envelope, DHTFindNode) {
		err = fmt.Errorf("calling from node %s is not permitted", req.NodeID)
		log.Error(err)
		return
	}
	node, err := DHT.Consistent.GetNode(string(req.NodeID))
	if err != nil {
		err = fmt.Errorf("get node %s from DHT failed: %s", req.NodeID, err)
		log.Error(err)
		return
	}
	resp.Node = node
	return
}

// FindNeighbor RPC returns FindNeighborReq.Count closest node from DHT
func (DHT *DHTService) FindNeighbor(req *proto.FindNeighborReq, resp *proto.FindNeighborResp) (err error) {
	if !IsPermitted(&req.Envelope, DHTFindNeighbor) {
		err = fmt.Errorf("calling from node %s is not permitted", req.NodeID)
		log.Error(err)
		return
	}

	nodes, err := DHT.Consistent.GetNeighbors(string(req.NodeID), req.Count)
	if err != nil {
		err = fmt.Errorf("get nodes from DHT failed: %s", err)
		log.Error(err)
		return
	}
	resp.Nodes = nodes
	return
}

// Ping RPC adds PingReq.Node to DHT
func (DHT *DHTService) Ping(req *proto.PingReq, resp *proto.PingResp) (err error) {
	log.Debugf("got req: %#v", req)
	if !IsPermitted(&req.Envelope, DHTPing) {
		err = fmt.Errorf("calling Ping from node %s is not permitted", req.NodeID)
		log.Error(err)
		return
	}

	// Checking if ID Nonce Pubkey matched
	if !kms.IsIDPubNonceValid(req.Node.ID.ToRawNodeID(), &req.Node.Nonce, req.Node.PublicKey) {
		err = fmt.Errorf("node: %s nonce public key not match", req.Node.ID)
		log.Error(err)
		return
	}

	// Checking MinNodeIDDifficulty
	if req.Node.ID.Difficulty() < conf.GConf.MinNodeIDDifficulty {
		err = fmt.Errorf("node: %s difficulty too low", req.Node.ID)
		log.Error(err)
		return
	}

	err = DHT.Consistent.Add(req.Node)
	if err != nil {
		err = fmt.Errorf("DHT.Consistent.Add %v failed: %s", req.Node, err)
	} else {
		resp.Msg = "Pong"
	}
	return
}
