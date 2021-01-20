// Copyright (c) 2020-2021 The bitcoinpay developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"github.com/btceasypay/bitcoinpay/core/types/pow"
	"github.com/btceasypay/bitcoinpay/crypto/seed"
	"github.com/btceasypay/bitcoinpay/params"
	"github.com/btceasypay/bitcoinpay/bx"
	"github.com/btceasypay/bitcoinpay/wallet"
	"io/ioutil"
	"os"
	"strings"
)

const (
	BX_VERSION = "0.0.1"
	TX_VERION  = 1 //default version is 1
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: bx [--version] [--help] <command> [<args>]\n")
	fmt.Fprintf(os.Stderr, `
encode and decode :
    base58-encode         encode a base16 string to a base58 string
    base58-decode         decode a base58 string to a base16 string
    base58check-encode    encode a base58check string
    base58check-decode    decode a base58check string
    base64-encode         encode a base16 string to a base64 string
    base64-decode         encode a base64 string to a base16 string
    rlp-encode            encode a string to a rlp encoded base16 string 
    rlp-decode            decode a rlp base16 string to a human-readble representation

hash :
    blake2b256            calculate Blake2b 256 hash of a base16 data.
    blake2b512            calculate Blake2b 512 hash of a base16 data.
    sha256                calculate SHA256 hash of a base16 data. 
    sha3-256              calculate SHA3 256 hash of a base16 data.
    keccak-256            calculate legacy keccak 256 hash of a bash16 data.
    blake256              calculate blake256 hash of a base16 data.
    ripemd160             calculate ripemd160 hash of a base16 data.
    bitcoin160            calculate ripemd160(sha256(data))   
    hash160               calculate ripemd160(blake2b256(data))

compact :
	compact-to-target      convert compact to target
	target-to-compact      convert target to compact
cuckoo-difficulty :
    compact-to-gps          convert cuckoo compact to gps.
    gps-to-compact          convert cuckoo gps to compact.

hash-difficulty :
    compact-to-hashrate          convert compact to hashrate.
    hashrate-to-compact          convert hashrate to compact.

entropy (seed) & mnemoic & hd & ec 
    entropy               generate a cryptographically secure pseudorandom entropy (seed)
    hd-new                create a new HD(BIP32) private key from an entropy (seed)
    hd-to-ec              convert the HD (BIP32) format private/public key to a EC private/public key
    hd-to-public          derive the HD (BIP32) public key from a HD private key
    hd-decode             decode a HD (BIP32) private/public key serialization format
    hd-derive             Derive a child HD (BIP32) key from another HD public or private key.
    mnemonic-new          create a mnemonic world-list (BIP39) from an entropy
    mnemonic-to-entropy   return back to the entropy (the random seed) from a mnemonic world list (BIP39)
    mnemonic-to-seed      convert a mnemonic world-list (BIP39) to its 512 bits seed 
    ec-new                create a new EC private key from an entropy (seed).
    ec-to-public          derive the EC public key from an EC private key (the compressed format by default )
    ec-to-wif             convert an EC private key to a WIF, associates with the compressed public key by default.
    wif-to-ec             convert a WIF private key to an EC private key.
    wif-to-public         derive the EC public key from a WIF private key. 

addr & tx & sign
    ec-to-addr            convert an EC public key to a paymant address. default is bx address
    tx-encode             encode a unsigned transaction.
    tx-decode             decode a transaction in base16 to json format.
    tx-sign               sign a transactions using a private key.
    msg-sign              create a message signature
    msg-verify            validate a message signature
    signature-decode      decode a ECDSA signature
	
`)
	os.Exit(1)
}

func cmdUsage(cmd *flag.FlagSet, usage string) {
	fmt.Fprintf(os.Stderr, usage)
	cmd.PrintDefaults()
}

func version() {
	fmt.Fprintf(os.Stderr, "bx Version : %q\n", BX_VERSION)
	os.Exit(1)
}

func errExit(err error) {
	fmt.Fprintf(os.Stderr, "bx Error : %q\n", err)
	os.Exit(1)
}

var base58checkVersion bx.BitcoinpayBase58checkVersionFlag
var base58checkVersionSize int
var base58checkMode string
var edgeBits int
var blocktime int
var unit string
var cuckoo_blocktime int
var mode string
var printDetail bool
var powtype string
var mheight int
var showDetails bool
var base58checkHasher string
var base58checkCksumSize int
var seedSize uint
var hdVer bx.Bip32VersionFlag
var hdHarden bool
var hdIndex uint
var derivePath bx.DerivePathFlag
var mnemoicSeedPassphrase string
var curve string
var uncompressedPKFormat bool
var network string
var powType string
var txInputs bx.TxInputsFlag
var txOutputs bx.TxOutputsFlag
var txVersion bx.TxVersionFlag
var txLockTime bx.TxLockTimeFlag
var privateKey string
var msgSignatureMode string

func main() {

	// ----------------------------
	// cmd for encoding & decoding
	// ----------------------------

	base58CheckEncodeCommand := flag.NewFlagSet("base58check-encode", flag.ExitOnError)
	base58checkVersion = bx.BitcoinpayBase58checkVersionFlag{}
	base58checkVersion.Set("testnet")
	base58CheckEncodeCommand.Var(&base58checkVersion, "v", "base58check `version` [mainnet|testnet|privnet")
	base58CheckEncodeCommand.StringVar(&base58checkMode, "m", "bitcoinpay", "base58check encode mode : [bitcoinpay|btc]")
	base58CheckEncodeCommand.StringVar(&base58checkHasher, "a", "", "base58check hasher")
	base58CheckEncodeCommand.IntVar(&base58checkCksumSize, "c", 4, "base58check checksum size")
	base58CheckEncodeCommand.Usage = func() {
		cmdUsage(base58CheckEncodeCommand, "Usage: bx base58check-encode [-v <ver>] [hexstring]\n")
	}

	base58CheckDecodeCommand := flag.NewFlagSet("base58check-decode", flag.ExitOnError)
	base58CheckDecodeCommand.BoolVar(&showDetails, "d", false, "show decode details")
	base58CheckDecodeCommand.StringVar(&base58checkMode, "m", "bitcoinpay", "base58check decode `mode`: [bitcoinpay|btc]")
	base58CheckDecodeCommand.StringVar(&base58checkHasher, "a", "", "base58check `hasher`")
	base58CheckDecodeCommand.IntVar(&base58checkVersionSize, "vs", 2, "base58check version `size`")
	base58CheckDecodeCommand.IntVar(&base58checkCksumSize, "cs", 4, "base58check checksum `size`")
	base58CheckDecodeCommand.Usage = func() {
		cmdUsage(base58CheckDecodeCommand, "Usage: bx base58check-decode [hexstring]\n")
	}

	base58EncodeCmd := flag.NewFlagSet("base58-encode", flag.ExitOnError)
	base58EncodeCmd.Usage = func() {
		cmdUsage(base58EncodeCmd, "Usage: bx base58-encode [hexstring]\n")
	}
	base58DecodeCmd := flag.NewFlagSet("base58-decode", flag.ExitOnError)
	base58DecodeCmd.Usage = func() {
		cmdUsage(base58DecodeCmd, "Usage: bx base58-decode [hexstring]\n")
	}

	compactToGPSCmd := flag.NewFlagSet("compact-to-gps", flag.ExitOnError)
	compactToGPSCmd.IntVar(&edgeBits, "e", 29, "edgebits")
	compactToGPSCmd.BoolVar(&printDetail, "p", false, "-p means print details")
	compactToGPSCmd.IntVar(&cuckoo_blocktime, "t", 43, "blocktime")
	compactToGPSCmd.IntVar(&mheight, "m", 1, "mheight")
	compactToGPSCmd.StringVar(&network, "n", "testnet", "the target network. (mainnet, testnet, privnet,mixnet)")
	compactToGPSCmd.StringVar(&powType, "a", "cuckaroom", "-a means use the selected cuckoo pow algorithm (cuckaroo,cuckaroom, cuckatoo), default is cuckaroom")
	compactToGPSCmd.Usage = func() {
		cmdUsage(compactToGPSCmd, "Usage: bx compact-to-gps -e 24 -ct 43 [compact value]\n")
	}

	gpsToCompactCmd := flag.NewFlagSet("gps-to-compact", flag.ExitOnError)
	gpsToCompactCmd.IntVar(&edgeBits, "e", 29, "edgebits")
	gpsToCompactCmd.IntVar(&cuckoo_blocktime, "t", 43, "blocktime")
	gpsToCompactCmd.IntVar(&mheight, "m", 1, "mheight")
	gpsToCompactCmd.StringVar(&network, "n", "testnet", "the target network. (mainnet, testnet, privnet,mixnet)")
	gpsToCompactCmd.StringVar(&powType, "a", "cuckaroom", "-a means use the selected cuckoo pow algorithm (cuckaroo,cuckaroom, cuckatoo), default is cuckaroom")
	gpsToCompactCmd.Usage = func() {
		cmdUsage(gpsToCompactCmd, "Usage: bx gps-to-compact -e 29 -t 43 [GPS float64]\n")
	}

	hashrateToCompactCmd := flag.NewFlagSet("hashrate-to-compact", flag.ExitOnError)
	hashrateToCompactCmd.IntVar(&blocktime, "t", 100, "blocktime")
	hashrateToCompactCmd.Usage = func() {
		cmdUsage(hashrateToCompactCmd, "Usage: bx hashrate-to-compact [Input unit should only be hash/s, no-integer input results a error]\n")
	}

	compactToHashrateCmd := flag.NewFlagSet("compact-to-hashrate", flag.ExitOnError)
	compactToHashrateCmd.BoolVar(&printDetail, "p", false, "-p means print details")
	compactToHashrateCmd.IntVar(&blocktime, "t", 100, "blocktime")
	compactToHashrateCmd.StringVar(&unit, "u", "H", "-u mens unit, (H, K, M, G, T,P)  H means Hash/s, it's by default")
	compactToHashrateCmd.Usage = func() {
		cmdUsage(compactToHashrateCmd, "Usage: bx compact-to-hashrate [compact]\n")
	}

	targetToCompactCmd := flag.NewFlagSet("target-to-compact", flag.ExitOnError)
	targetToCompactCmd.StringVar(&mode, "pow", "hash", "pow type (hash , cuckoo24,cuckoo29)")
	targetToCompactCmd.Usage = func() {
		cmdUsage(targetToCompactCmd, "Usage: bx target-to-compact [target]\n")
	}

	compactToTargetCmd := flag.NewFlagSet("compact-to-target", flag.ExitOnError)
	compactToTargetCmd.StringVar(&powtype, "pow", "hash", "pow type (hash , cuckoo24,cuckoo29)")
	compactToTargetCmd.Usage = func() {
		cmdUsage(compactToTargetCmd, "Usage: bx compact-to-target [compact]\n")
	}

	base64EncodeCmd := flag.NewFlagSet("base64-encode", flag.ExitOnError)
	base64EncodeCmd.Usage = func() {
		cmdUsage(base64EncodeCmd, "Usage: bx base64-encode [hexstring]\n")
	}
	base64DecodeCmd := flag.NewFlagSet("base64-decode", flag.ExitOnError)
	base64DecodeCmd.Usage = func() {
		cmdUsage(base64DecodeCmd, "Usage: bx base64-decode [hexstring]\n")
	}

	rlpEncodeCmd := flag.NewFlagSet("rlp-encode", flag.ExitOnError)
	rlpEncodeCmd.Usage = func() {
		cmdUsage(rlpEncodeCmd, "Usage: bx rlp-encode [string]\n")
	}

	rlpDecodeCmd := flag.NewFlagSet("rlp-decode", flag.ExitOnError)
	rlpDecodeCmd.Usage = func() {
		cmdUsage(rlpDecodeCmd, "Usage: bx rlp-decode [hexstring]\n")
	}

	// ----------------------------
	// cmd for hashing
	// ----------------------------

	sha256cmd := flag.NewFlagSet("sha256", flag.ExitOnError)
	sha256cmd.Usage = func() {
		cmdUsage(sha256cmd, "Usage: bx sha256 [hexstring]\n")
	}

	blake2b256cmd := flag.NewFlagSet("blake2b256", flag.ExitOnError)
	blake2b256cmd.Usage = func() {
		cmdUsage(blake2b256cmd, "Usage: bx blak2b256 [hexstring]\n")
	}

	blake2b512cmd := flag.NewFlagSet("blake2b512", flag.ExitOnError)
	blake2b512cmd.Usage = func() {
		cmdUsage(blake2b512cmd, "Usage: bx blak2b512 [hexstring]\n")
	}

	blake256cmd := flag.NewFlagSet("blake256", flag.ExitOnError)
	blake256cmd.Usage = func() {
		cmdUsage(blake256cmd, "Usage: bx blake256 [hexstring]\n")
	}

	sha3_256cmd := flag.NewFlagSet("sha3-256", flag.ExitOnError)
	sha3_256cmd.Usage = func() {
		cmdUsage(sha3_256cmd, "Usage: bx sha3-256 [hexstring]\n")
	}

	keccak256cmd := flag.NewFlagSet("keccak-256", flag.ExitOnError)
	keccak256cmd.Usage = func() {
		cmdUsage(keccak256cmd, "Usage: bx keccak-256 [hexstring]\n")
	}

	ripemd160Cmd := flag.NewFlagSet("ripemd160", flag.ExitOnError)
	ripemd160Cmd.Usage = func() {
		cmdUsage(ripemd160Cmd, "Usage: bx ripemd160 [hexstring]\n")
	}

	bitcion160Cmd := flag.NewFlagSet("bitcoin160", flag.ExitOnError)
	bitcion160Cmd.Usage = func() {
		cmdUsage(bitcion160Cmd, "Usage: bx bitcoin160 [hexstring]\n")
	}

	hash160Cmd := flag.NewFlagSet("hash160", flag.ExitOnError)
	hash160Cmd.Usage = func() {
		cmdUsage(bitcion160Cmd, "Usage: bx hash160 [hexstring]\n")
	}

	// ----------------------------
	// cmd for crypto
	// ----------------------------

	// Entropy (Seed)
	entropyCmd := flag.NewFlagSet("entropy", flag.ExitOnError)
	entropyCmd.Usage = func() {
		cmdUsage(entropyCmd, "Usage: bx entropy [-s size] \n")
	}
	entropyCmd.UintVar(&seedSize, "s", seed.DefaultSeedBytes*8, "The length in bits for a seed (entropy)")

	// HD (BIP32)
	hdNewCmd := flag.NewFlagSet("hd-new", flag.ExitOnError)
	hdNewCmd.Usage = func() {
		cmdUsage(hdNewCmd, "Usage: bx hd-new [-v version] [entropy] \n")
	}
	hdVer.Set("testnet")
	hdNewCmd.Var(&hdVer, "v", "The HD(BIP32) `version` [mainnet|testnet|privnet|bip32]")

	hdToPubCmd := flag.NewFlagSet("hd-to-public", flag.ExitOnError)
	hdToPubCmd.Usage = func() {
		cmdUsage(hdToPubCmd, "Usage: bx hd-to-public [hd_private_key] \n")
	}
	hdToPubCmd.Var(&hdVer, "v", "The HD(BIP32) `version` [mainnet|testnet|privnet|bip32]")

	hdToEcCmd := flag.NewFlagSet("hd-to-ec", flag.ExitOnError)
	hdToEcCmd.Usage = func() {
		cmdUsage(hdToEcCmd, "Usage: bx hd-to-ec [hd_private_key or hd_public_key] \n")
	}
	hdToEcCmd.Var(&hdVer, "v", "The HD(BIP32) `version` [mainnet|testnet|privnet|bip32]")

	hdDecodeCmd := flag.NewFlagSet("hd-decode", flag.ExitOnError)
	hdDecodeCmd.Usage = func() {
		cmdUsage(hdDecodeCmd, "Usage: bx hd-decode [hd_private_key or hd_public_key] \n")
	}

	hdDeriveCmd := flag.NewFlagSet("hd-derive", flag.ExitOnError)
	hdDeriveCmd.Usage = func() {
		cmdUsage(hdDeriveCmd, "Usage: bx hd-derive [hd_private_key or hd_public_key] \n")
	}
	hdDeriveCmd.UintVar(&hdIndex, "i", 0, "The HD `index`")
	hdDeriveCmd.BoolVar(&hdHarden, "d", false, "create a hardened key")
	derivePath = bx.DerivePathFlag{Path: wallet.DerivationPath{}}
	hdDeriveCmd.Var(&derivePath, "p", "hd derive `path`. ex: m/44'/0'/0'/0")
	hdDeriveCmd.Var(&hdVer, "v", "The HD(BIP32) `version` [mainnet|testnet|privnet|bip32]")

	// Mnemonic (BIP39)
	mnemonicNewCmd := flag.NewFlagSet("mnemonic-new", flag.ExitOnError)
	mnemonicNewCmd.Usage = func() {
		cmdUsage(mnemonicNewCmd, "Usage: bx mnemonic-new [entropy]  \n")
	}

	mnemonicToEntropyCmd := flag.NewFlagSet("mnemonic-to-entropy", flag.ExitOnError)
	mnemonicToEntropyCmd.Usage = func() {
		cmdUsage(mnemonicToEntropyCmd, "Usage: bx mnemonic-to-entropy [mnemonic]  \n")
	}

	mnemonicToSeedCmd := flag.NewFlagSet("mnemonic-to-seed", flag.ExitOnError)
	mnemonicToSeedCmd.Usage = func() {
		cmdUsage(mnemonicToSeedCmd, "Usage: bx mnemonic-to-seed [mnemonic]  \n")
	}
	mnemonicToSeedCmd.StringVar(&mnemoicSeedPassphrase, "p", "", "An optional passphrase for converting the mnemonic to a seed")

	// EC
	ecNewCmd := flag.NewFlagSet("ec-new", flag.ExitOnError)
	ecNewCmd.Usage = func() {
		cmdUsage(ecNewCmd, "Usage: bx ec-new [entropy]  \n")
	}
	ecNewCmd.StringVar(&curve, "c", "secp256k1", "the elliptic curve is using")

	ecToPubCmd := flag.NewFlagSet("ec-to-public", flag.ExitOnError)
	ecToPubCmd.Usage = func() {
		cmdUsage(ecToPubCmd, "Usage: bx ec-to-public [ec_private_key] \n")
	}
	ecToPubCmd.BoolVar(&uncompressedPKFormat, "u", false, "using the uncompressed public key format")

	// Wif
	ecToWifCmd := flag.NewFlagSet("ec-to-wif", flag.ExitOnError)
	ecToWifCmd.Usage = func() {
		cmdUsage(ecToWifCmd, "Usage: bx ec-to-wif [ec_private_key] \n")
	}
	ecToWifCmd.BoolVar(&uncompressedPKFormat, "u", false, "using the uncompressed public key format")

	wifToEcCmd := flag.NewFlagSet("wif-to-ec", flag.ExitOnError)
	wifToEcCmd.Usage = func() {
		cmdUsage(wifToEcCmd, "Usage: bx wif-to-ec [WIF] \n")
	}

	wifToPubCmd := flag.NewFlagSet("wif-to-public", flag.ExitOnError)
	wifToPubCmd.Usage = func() {
		cmdUsage(wifToPubCmd, "Usage: bx wif-to-public [WIF] \n")
	}
	wifToPubCmd.BoolVar(&uncompressedPKFormat, "u", false, "using the uncompressed public key format")

	// Address
	ecToAddrCmd := flag.NewFlagSet("ec-to-addr", flag.ExitOnError)
	ecToAddrCmd.Usage = func() {
		cmdUsage(ecToAddrCmd, "Usage: bx ec-to-addr [ec_public_key] \n")
	}
	ecToAddrCmd.Var(&base58checkVersion, "v", "base58check `version` [mainnet|testnet|privnet]")

	// Transaction
	txDecodeCmd := flag.NewFlagSet("tx-decode", flag.ExitOnError)
	txDecodeCmd.Usage = func() {
		cmdUsage(txDecodeCmd, "Usage: bx tx-decode [base16_string] \n")
	}
	txDecodeCmd.StringVar(&network, "n", "testnet", "decode rawtx for the target network. (mainnet, testnet, privnet)")

	txEncodeCmd := flag.NewFlagSet("tx-encode", flag.ExitOnError)
	txEncodeCmd.Usage = func() {
		cmdUsage(txEncodeCmd, "Usage: bx tx-encode [-i tx-input] [-l tx-lock-time] [-o tx-output] [-v tx-version] \n")
	}
	txVersion = bx.TxVersionFlag(TX_VERION) //set default tx version
	txEncodeCmd.Var(&txVersion, "v", "the transaction version")
	txEncodeCmd.Var(&txLockTime, "l", "the transaction lock time")
	txEncodeCmd.Var(&txInputs, "i", `The set of transaction input points encoded as TXHASH:INDEX:SEQUENCE. 
TXHASH is a Base16 transaction hash. INDEX is the 32 bit input index
in the context of the transaction. SEQUENCE is the optional 32 bit 
input sequence and defaults to the maximum value.`)
	txEncodeCmd.Var(&txOutputs, "o", `The set of transaction output data encoded as TARGET:BPAY. 
TARGET is an address (pay-to-pubkey-hash or pay-to-script-hash).
BPAY is the 64 bit spend amount in bitcoinpay.`)

	txSignCmd := flag.NewFlagSet("tx-sign", flag.ExitOnError)
	txSignCmd.Usage = func() {
		cmdUsage(txSignCmd, "Usage: bx tx-sign [raw_tx_base16_string] \n")
	}
	txSignCmd.StringVar(&privateKey, "k", "", "the ec private key to sign the raw transaction")

	msgSignCmd := flag.NewFlagSet("msg-sign", flag.ExitOnError)
	msgSignCmd.Usage = func() {
		cmdUsage(msgSignCmd, "Usage: msg-sign [wif] [message] \n")
	}
	msgSignCmd.StringVar(&msgSignatureMode, "m", "bx", "the msg signature mode")
	msgSignCmd.BoolVar(&showDetails, "d", false, "show signature details")

	msgVerifyCmd := flag.NewFlagSet("msg-verify", flag.ExitOnError)
	msgVerifyCmd.Usage = func() {
		cmdUsage(msgVerifyCmd, "Usage: msg-verify [addr] [signature] [message] \n")
	}
	msgVerifyCmd.StringVar(&msgSignatureMode, "m", "bx", "the msg signature mode")

	flagSet := []*flag.FlagSet{
		base58CheckEncodeCommand,
		base58CheckDecodeCommand,
		base58EncodeCmd,
		base58DecodeCmd,
		gpsToCompactCmd,
		compactToGPSCmd,
		compactToTargetCmd,
		targetToCompactCmd,
		compactToHashrateCmd,
		hashrateToCompactCmd,
		base64EncodeCmd,
		base64DecodeCmd,
		rlpEncodeCmd,
		rlpDecodeCmd,
		sha256cmd,
		blake2b256cmd,
		blake2b512cmd,
		blake256cmd,
		sha3_256cmd,
		keccak256cmd,
		ripemd160Cmd,
		bitcion160Cmd,
		hash160Cmd,
		entropyCmd,
		hdNewCmd,
		hdToPubCmd,
		hdToEcCmd,
		hdDecodeCmd,
		hdDeriveCmd,
		mnemonicNewCmd,
		mnemonicToEntropyCmd,
		mnemonicToSeedCmd,
		ecNewCmd,
		ecToPubCmd,
		ecToWifCmd,
		wifToEcCmd,
		wifToPubCmd,
		ecToAddrCmd,
		txEncodeCmd,
		txDecodeCmd,
		txSignCmd,
		msgSignCmd,
		msgVerifyCmd,
	}

	if len(os.Args) == 1 {
		usage()
	}
	switch os.Args[1] {
	case "help", "--help":
		usage()
	case "version", "--version":
		version()
	default:
		valid := false
		for _, cmd := range flagSet {
			if os.Args[1] == cmd.Name() {
				cmd.Parse(os.Args[2:])
				valid = true
				break
			}
		}
		if !valid {
			invalid := os.Args[1]
			if invalid[0] == '-' {
				fmt.Fprintf(os.Stderr, "unknown option: %q \n", invalid)
			} else {
				fmt.Fprintf(os.Stderr, "%q is not valid command\n", invalid)
			}
			os.Exit(1)
		}
	}
	// Handle base58check-encode
	if base58CheckEncodeCommand.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base58CheckEncodeCommand.Usage()
			} else {
				bx.Base58CheckEncode(base58checkVersion.Ver, base58checkMode, base58checkHasher, base58checkCksumSize, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base58CheckEncode(base58checkVersion.Ver, base58checkMode, base58checkHasher, base58checkCksumSize, str)
		}
	}

	// Handle base58check-decode
	if base58CheckDecodeCommand.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base58CheckDecodeCommand.Usage()
			} else {
				bx.Base58CheckDecode(base58checkMode, base58checkHasher, base58checkVersionSize, base58checkCksumSize, os.Args[len(os.Args)-1], showDetails)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base58CheckDecode(base58checkMode, base58checkHasher, base58checkVersionSize, base58checkCksumSize, str, showDetails)
		}
	}

	// Handle base58-encode
	if base58EncodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base58EncodeCmd.Usage()
			} else {
				bx.Base58Encode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base58Encode(str)
		}
	}

	getNetWork := func(network string) *params.Params {
		switch network {
		case "testnet":
			return &params.TestNetParams
		case "privnet":
			return &params.PrivNetParams
		case "mainnet":
			return &params.MainNetParams
		case "mixnet":
			return &params.MixNetParams
		default:
			return &params.TestNetParams
		}
	}
	getCuckooScale := func(powType string, p *params.Params, edgeBits, mheight int64) int {
		switch powType {
		case "cuckaroo":
			instance := &pow.Cuckaroo{}
			instance.SetMainHeight(mheight)
			instance.SetEdgeBits(uint8(edgeBits))
			instance.SetParams(p.PowConfig)
			return int(instance.GraphWeight())
		case "cuckaroom":
			instance := &pow.Cuckaroom{}
			instance.SetMainHeight(mheight)
			instance.SetEdgeBits(uint8(edgeBits))
			instance.SetParams(p.PowConfig)
			return int(instance.GraphWeight())
		case "cuckatoo":
			instance := &pow.Cuckaroo{}
			instance.SetMainHeight(mheight)
			instance.SetEdgeBits(uint8(edgeBits))
			instance.SetParams(p.PowConfig)
			return int(instance.GraphWeight())
		}
		return 0
	}
	// Handle compact-to-gps
	if compactToGPSCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		p := getNetWork(network)
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				compactToGPSCmd.Usage()
			} else {
				bx.CompactToGPS(os.Args[len(os.Args)-1], cuckoo_blocktime, getCuckooScale(powType, p, int64(edgeBits), int64(mheight)), printDetail)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.CompactToGPS(str, cuckoo_blocktime, getCuckooScale(powType, p, int64(edgeBits), int64(mheight)), printDetail)
		}
	}
	// Handle gps-to-compact
	if gpsToCompactCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		p := getNetWork(network)
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				gpsToCompactCmd.Usage()
			} else {
				bx.GPSToCompact(os.Args[len(os.Args)-1], cuckoo_blocktime, getCuckooScale(powType, p, int64(edgeBits), int64(mheight)))
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.GPSToCompact(str, cuckoo_blocktime, getCuckooScale(powType, p, int64(edgeBits), int64(mheight)))
		}
	}

	// Handle hashrate- to compact-
	if hashrateToCompactCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hashrateToCompactCmd.Usage()
			} else {
				bx.HashrateToCompact(os.Args[len(os.Args)-1], blocktime)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HashrateToCompact(str, blocktime)
		}
	}
	// Handle hash compact- to difficulty
	if compactToHashrateCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				compactToHashrateCmd.Usage()
			} else {
				bx.CompactToHashrate(os.Args[len(os.Args)-1], unit, printDetail, blocktime)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.CompactToHashrate(str, unit, printDetail, blocktime)
		}
	}
	// Handle hash compact- to target
	if compactToTargetCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				compactToTargetCmd.Usage()
			} else {
				bx.CompactToTarget(os.Args[len(os.Args)-1], powtype)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.CompactToTarget(str, powtype)
		}
	}
	// Handle hash target- to compact
	if targetToCompactCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				targetToCompactCmd.Usage()
			} else {
				bx.TargetToCompact(os.Args[len(os.Args)-1], powtype)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.TargetToCompact(str, powtype)
		}
	}

	// Handle base58-decode
	if base58DecodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base58DecodeCmd.Usage()
			} else {
				bx.Base58Decode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base58Decode(str)
		}
	}

	if base64EncodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base64EncodeCmd.Usage()
			} else {
				bx.Base64Encode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base64Encode(str)
		}
	}

	if base64DecodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				base64DecodeCmd.Usage()
			} else {
				bx.Base64Decode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Base64Decode(str)
		}
	}

	if rlpEncodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				rlpEncodeCmd.Usage()
			} else {
				bx.RlpEncode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.RlpEncode(str)
		}
	}
	if rlpDecodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				rlpDecodeCmd.Usage()
			} else {
				bx.RlpDecode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.RlpDecode(str)
		}
	}

	if sha256cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				sha256cmd.Usage()
			} else {
				bx.Sha256(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Sha256(str)
		}
	}

	if blake256cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				blake256cmd.Usage()
			} else {
				bx.Blake256(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Blake256(str)
		}
	}

	if blake2b256cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				blake2b256cmd.Usage()
			} else {
				bx.Blake2b256(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Blake2b256(str)
		}
	}

	if blake2b512cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				blake2b512cmd.Usage()
			} else {
				bx.Blake2b512(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Blake2b512(str)
		}
	}

	if sha3_256cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				sha3_256cmd.Usage()
			} else {
				bx.Sha3_256(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Sha3_256(str)
		}
	}

	if keccak256cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				keccak256cmd.Usage()
			} else {
				bx.Keccak256(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Keccak256(str)
		}
	}

	if ripemd160Cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				ripemd160Cmd.Usage()
			} else {
				bx.Ripemd160(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Ripemd160(str)
		}
	}

	if bitcion160Cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				bitcion160Cmd.Usage()
			} else {
				bx.Bitcoin160(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Bitcoin160(str)
		}
	}

	if hash160Cmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hash160Cmd.Usage()
			} else {
				bx.Hash160(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.Hash160(str)
		}
	}

	if entropyCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) > 2 && (os.Args[2] == "help" || os.Args[2] == "--help") {
				entropyCmd.Usage()
			} else {
				if seedSize%8 > 0 {
					errExit(fmt.Errorf("seed (entropy) length must be Must be divisible by 8"))
				}
				s, err := bx.NewEntropy(seedSize / 8)
				if err != nil {
					bx.ErrExit(err)
				}
				fmt.Printf("%s\n", s)
			}
		} else {
			entropyCmd.Usage()
		}
	}

	if hdNewCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hdNewCmd.Usage()
			} else {
				bx.HdNewMasterPrivateKey(hdVer.Version, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HdNewMasterPrivateKey(hdVer.Version, str)
		}
	}

	if hdToPubCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hdToPubCmd.Usage()
			} else {
				bx.HdPrivateKeyToHdPublicKey(hdVer.Version, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HdPrivateKeyToHdPublicKey(hdVer.Version, str)
		}
	}

	if hdToEcCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hdToEcCmd.Usage()
			} else {
				bx.HdKeyToEcKey(hdVer.Version, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HdKeyToEcKey(hdVer.Version, str)
		}
	}

	if hdDecodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hdDecodeCmd.Usage()
			} else {
				bx.HdDecode(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HdDecode(str)
		}
	}

	if hdDeriveCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				hdDeriveCmd.Usage()
			} else {
				bx.HdDerive(hdHarden, uint32(hdIndex), derivePath.Path, hdVer.Version, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.HdDerive(hdHarden, uint32(hdIndex), derivePath.Path, hdVer.Version, str)
		}
	}

	if mnemonicNewCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				mnemonicNewCmd.Usage()
			} else {
				bx.MnemonicNew(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.MnemonicNew(str)
		}
	}

	if mnemonicToEntropyCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				mnemonicToEntropyCmd.Usage()
			} else {
				bx.MnemonicToEntropy(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.MnemonicToEntropy(str)
		}
	}

	if mnemonicToSeedCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				mnemonicToSeedCmd.Usage()
			} else {
				bx.MnemonicToSeed(mnemoicSeedPassphrase, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.MnemonicToSeed(mnemoicSeedPassphrase, str)
		}
	}

	if ecNewCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				ecNewCmd.Usage()
			} else {
				pk, err := bx.EcNew(curve, os.Args[len(os.Args)-1])
				if err != nil {
					bx.ErrExit(err)
				}
				fmt.Printf("%s\n", pk)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			pk, err := bx.EcNew(curve, str)
			if err != nil {
				bx.ErrExit(err)
			}
			fmt.Printf("%s\n", pk)
		}
	}

	if ecToPubCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				ecToPubCmd.Usage()
			} else {
				key, err := bx.EcPrivateKeyToEcPublicKey(uncompressedPKFormat, os.Args[len(os.Args)-1])
				if err != nil {
					bx.ErrExit(err)
				}
				fmt.Printf("%s\n", key)

			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			key, err := bx.EcPrivateKeyToEcPublicKey(uncompressedPKFormat, str)
			if err != nil {
				bx.ErrExit(err)
			}
			fmt.Printf("%s\n", key)
		}
	}

	if ecToWifCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				ecToWifCmd.Usage()
			} else {
				bx.EcPrivateKeyToWif(uncompressedPKFormat, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.EcPrivateKeyToWif(uncompressedPKFormat, str)
		}
	}

	if wifToEcCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				wifToEcCmd.Usage()
			} else {
				bx.WifToEcPrivateKey(os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.WifToEcPrivateKey(str)
		}
	}

	if wifToPubCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				wifToPubCmd.Usage()
			} else {
				bx.WifToEcPubkey(uncompressedPKFormat, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.WifToEcPubkey(uncompressedPKFormat, str)
		}
	}

	if ecToAddrCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				ecToAddrCmd.Usage()
			} else {
				bx.EcPubKeyToAddressSTDO(base58checkVersion.Ver, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.EcPubKeyToAddressSTDO(base58checkVersion.Ver, str)
		}
	}

	if txDecodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				txDecodeCmd.Usage()
			} else {
				bx.TxDecode(network, os.Args[len(os.Args)-1])
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.TxDecode(network, str)
		}
	}

	if txEncodeCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				txEncodeCmd.Usage()
			} else {
				bx.TxEncodeSTDO(txVersion, txLockTime, txInputs, txOutputs)
			}
		}
	}

	if txSignCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				txSignCmd.Usage()
			} else {
				bx.TxSignSTDO(privateKey, os.Args[len(os.Args)-1], network)
			}
		} else { //try from STDIN
			src, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				errExit(err)
			}
			str := strings.TrimSpace(string(src))
			bx.TxSignSTDO(privateKey, str, network)
		}
	}

	if msgSignCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				msgSignCmd.Usage()
			} else {
				bx.MsgSign(msgSignatureMode, showDetails, os.Args[len(os.Args)-2], os.Args[len(os.Args)-1], showDetails)
			}
		}
	}

	if msgVerifyCmd.Parsed() {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			if len(os.Args) == 2 || os.Args[2] == "help" || os.Args[2] == "--help" {
				msgVerifyCmd.Usage()
			} else {
				bx.VerifyMsgSignature(msgSignatureMode, os.Args[len(os.Args)-3], os.Args[len(os.Args)-2], os.Args[len(os.Args)-1])
			}
		}
	}
}