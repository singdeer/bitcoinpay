// Copyright (c) 2020-2021 The bitcoinpay developers
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// The parts code inspired & originated from
// https://github.com/ethereum/go-ethereum/trie

package trie

import (
	"fmt"
	"github.com/btceasypay/bitcoinpay/common/hash"
)

// MissingNodeError is returned by the trie functions (TryGet, TryUpdate, TryDelete)
// in the case where a trie node is not present in the local database. It contains
// information necessary for retrieving the missing node.
type MissingNodeError struct {
	NodeHash hash.Hash // hash of the missing node
	Path     []byte    // hex-encoded path to the missing node
}

func (err *MissingNodeError) Error() string {
	return fmt.Sprintf("missing trie node %x (path %x)", err.NodeHash, err.Path)
}
