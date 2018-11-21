/*
 * Copyright 2018 The CovenantSQL Authors.
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

package types

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/CovenantSQL/CovenantSQL/crypto/asymmetric"
	"github.com/CovenantSQL/CovenantSQL/crypto/hash"
	"github.com/CovenantSQL/CovenantSQL/crypto/verifier"
	"github.com/CovenantSQL/CovenantSQL/utils"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSignAndVerify(t *testing.T) {
	block, err := createRandomBlock(genesisHash, true)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if err = block.Verify(); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	block.SignedHeader.HSV.DataHash[0]++

	if err = errors.Cause(block.Verify()); err != verifier.ErrHashValueNotMatch {
		t.Fatalf("Unexpected error: %v", err)
	}

	block.Acks = append(block.Acks, &SignedAckHeader{
		DefaultHashSignVerifierImpl: verifier.DefaultHashSignVerifierImpl{
			DataHash: hash.Hash{0x01},
		},
	})

	if err = block.Verify(); err != ErrMerkleRootVerification {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestHeaderMarshalUnmarshaler(t *testing.T) {
	block, err := createRandomBlock(genesisHash, false)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	origin := &block.SignedHeader.Header
	enc, err := utils.EncodeMsgPack(origin)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	dec := &Header{}
	if err = utils.DecodeMsgPack(enc.Bytes(), dec); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts1, err := origin.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts2, err := dec.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}

	if !reflect.DeepEqual(origin, dec) {
		t.Fatalf("Values don't match:\n\tv1 = %+v\n\tv2 = %+v", origin, dec)
	}
}

func TestSignedHeaderMarshaleUnmarshaler(t *testing.T) {
	block, err := createRandomBlock(genesisHash, true)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	origin := &block.SignedHeader
	enc, err := utils.EncodeMsgPack(origin)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	dec := &SignedHeader{}

	if err = utils.DecodeMsgPack(enc.Bytes(), dec); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts1, err := origin.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts2, err := dec.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}

	if !reflect.DeepEqual(origin.Header, dec.Header) {
		t.Fatalf("Values don't match:\n\tv1 = %+v\n\tv2 = %+v", origin.Header, dec.Header)
	}

	if err = origin.Verify(); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if err = dec.Verify(); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}
}

func TestBlockMarshalUnmarshaler(t *testing.T) {
	origin, err := createRandomBlock(genesisHash, false)
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}
	origin2, err := createRandomBlock(genesisHash, false)
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	blocks := make(Blocks, 0, 2)
	blocks = append(blocks, origin)
	blocks = append(blocks, origin2)
	blocks = append(blocks, nil)

	blocks2 := make(Blocks, 0, 2)
	blocks2 = append(blocks2, origin)
	blocks2 = append(blocks2, origin2)
	blocks2 = append(blocks2, nil)

	bts1, err := blocks.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts2, err := blocks2.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}

	enc, err := utils.EncodeMsgPack(origin)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	dec := &Block{}

	if err = utils.DecodeMsgPack(enc.Bytes(), dec); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts1, err = origin.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	bts2, err = dec.MarshalHash()
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}

	if !reflect.DeepEqual(origin, dec) {
		t.Fatalf("Values don't match:\n\tv1 = %+v\n\tv2 = %+v", origin, dec)
	}
}

func TestGenesis(t *testing.T) {
	genesis, err := createRandomBlock(genesisHash, true)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if err = genesis.VerifyAsGenesis(); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if err = genesis.SignedHeader.VerifyAsGenesis(); err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	// Test non-genesis block
	genesis, err = createRandomBlock(genesisHash, false)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	if err = genesis.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatal("Unexpected result: returned nil while expecting an error")
	}

	if err = genesis.SignedHeader.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatal("Unexpected result: returned nil while expecting an error")
	}

	// Test altered public key block
	genesis, err = createRandomBlock(genesisHash, true)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	_, pub, err := asymmetric.GenSecp256k1KeyPair()

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	genesis.SignedHeader.HSV.Signee = pub

	if err = genesis.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatal("Unexpected result: returned nil while expecting an error")
	}

	if err = genesis.SignedHeader.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatal("Unexpected result: returned nil while expecting an error")
	}

	// Test altered signature
	genesis, err = createRandomBlock(genesisHash, true)

	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}

	genesis.SignedHeader.HSV.Signature.R.Add(genesis.SignedHeader.HSV.Signature.R, big.NewInt(int64(1)))
	genesis.SignedHeader.HSV.Signature.S.Add(genesis.SignedHeader.HSV.Signature.S, big.NewInt(int64(1)))

	if err = genesis.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatalf("Unexpected error: %v", err)
	}

	if err = genesis.SignedHeader.VerifyAsGenesis(); err != nil {
		t.Logf("Error occurred as expected: %v", err)
	} else {
		t.Fatal("Unexpected result: returned nil while expecting an error")
	}
}

func Test(t *testing.T) {
	Convey("CalcNextID should return correct id of each testing block", t, func() {
		var (
			nextid uint64
			ok     bool

			cases = [...]struct {
				block  *Block
				nextid uint64
				ok     bool
			}{
				{
					block: &Block{
						QueryTxs: []*QueryAsTx{},
					},
					nextid: 0,
					ok:     false,
				}, {
					block: &Block{
						QueryTxs: nil,
					},
					nextid: 0,
					ok:     false,
				}, {
					block: &Block{
						QueryTxs: []*QueryAsTx{
							&QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: ReadQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 10),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 0,
									},
								},
							},
						},
					},
					nextid: 0,
					ok:     false,
				}, {
					block: &Block{
						QueryTxs: []*QueryAsTx{
							&QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: WriteQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 10),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 0,
									},
								},
							},
						},
					},
					nextid: 10,
					ok:     true,
				}, {
					block: &Block{
						QueryTxs: []*QueryAsTx{
							&QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: ReadQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 10),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 0,
									},
								},
							}, &QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: WriteQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 10),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 0,
									},
								},
							}, &QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: ReadQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 10),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 10,
									},
								},
							}, &QueryAsTx{
								Request: &Request{
									Header: SignedRequestHeader{
										RequestHeader: RequestHeader{
											QueryType: WriteQuery,
										},
									},
									Payload: RequestPayload{
										Queries: make([]Query, 20),
									},
								},
								Response: &SignedResponseHeader{
									ResponseHeader: ResponseHeader{
										LogOffset: 10,
									},
								},
							},
						},
					},
					nextid: 30,
					ok:     true,
				},
			}
		)

		for _, v := range cases {
			nextid, ok = v.block.CalcNextID()
			So(ok, ShouldEqual, v.ok)
			if ok {
				So(nextid, ShouldEqual, v.nextid)
			}
		}
	})
}