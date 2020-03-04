package btc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/multi-chain/account"
	"github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/core/states"
	"github.com/ontio/multi-chain/core/store/leveldbstore"
	"github.com/ontio/multi-chain/core/store/overlaydb"
	"github.com/ontio/multi-chain/native"
	ccmcom "github.com/ontio/multi-chain/native/service/cross_chain_manager/common"
	"github.com/ontio/multi-chain/native/service/governance/side_chain_manager"
	"github.com/ontio/multi-chain/native/service/header_sync/btc"
	hscom "github.com/ontio/multi-chain/native/service/header_sync/common"
	"github.com/ontio/multi-chain/native/service/utils"
	"github.com/ontio/multi-chain/native/storage"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	acct *account.Account = account.NewAccount("")

	rdm               = "552102dec9a415b6384ec0a9331d0cdf02020f0f1e5731c327b86e2b5a92455a289748210365b1066bcfa21987c3e207b92e309b95ca6bee5f1133cf04d6ed4ed265eafdbc21031104e387cd1a103c27fdc8a52d5c68dec25ddfb2f574fbdca405edfd8c5187de21031fdb4b44a9f20883aff505009ebc18702774c105cb04b1eecebcb294d404b1cb210387cda955196cc2b2fc0adbbbac1776f8de77b563c6d2a06a77d96457dc3d0d1f2102dd7767b6a7cc83693343ba721e0f5f4c7b4b8d85eeb7aec20d227625ec0f59d321034ad129efdab75061e8d4def08f5911495af2dae6d3e9a4b6e7aeb5186fa432fc57ae"
	fromBtcTxid       = "2587a59e8069c563d32de9d4a2b946760d740b6963566dd7b32d8ec549f2d238"
	fromBtcRawTx      = "010000000147d9b1bc6a52099f746863722282e3febc9ad3ad6b2eac0f2df6d2badf1df28a020000006b483045022100a1e573ba3589217e1b20d6ed53e2dda705deb3d284122c61987266e66aff074802200165734cf4519b560d806d392f10cec2aeb3071cf72c759a5abc9c33cd2f983f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff031027000000000000220020216a09cb8ee51da1a91ea8942552d7936c886a10b507299003661816c0e9f18b00000000000000003d6a3b6602000000000000000000000000000000149702640a6b971ca18efc20ad73ca4e8ba390c910145cd3143f91a13fe971043e1e4605c1c23b46bf44620e0700000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000"
	fromBtcProof      = "00000020775635e1ada1581f0fa6eff86bfc4720253c9c4fcd7165843e902600000000003faec6ef7165e988b344b553b15dff0d66eb62e71b1d93462c64b0eab1086852fef54c5effff001d6c2edd4f4d0200000b3a4a5328d2e6b72f26fb5f3aa6db80e8301c2746c5ce6e21813e884c3a08e96a8ba1ccfe764700b7d956acff0697680b0e9412517972d8e8a10c9ac37c96fd0c81a705037d9f8caaa679075d525cd12bbb698e6f6917e61aecbe3d529f65c7bd38d2f249c58e2db3d76d5663690b740d7646b9a2d4e92dd363c569809ea58725a47e736292f4a96de7c46462c53b823c1732cb2d863c402bc3dd96527e69305e0af15f3487c9093c59d7dc0e7fcde6db50354e73f640987e3305e917aad7531abaa9513d16228fb2b17c3cd04f9ec97e3c38de9dac7ff2af93184c338e86e6c2d8731bc8430a7f31bc050d11776d6e3b665951af070fe889cba7aec895e40b3e67107e62b1ec0ebef9a226abc458b55920060f0a5c06edd26432a987d3f6cfafefee3301b3281270ceb45e5831e435fa70056cd28927251eebe875f2fd810aa501cca43940fe1ba0e8be004e8f05f740b66e2e9a1a24b76bdd4c0c6d53230c1903ef2d00"
	fromBtcMerkleRoot = "526808b1eab0642c46931d1be762eb660dff5db153b544b388e96571efc6ae3f"
	toOntAddr         = "AdzZ2VKufdJWeB8t9a8biXoHbbMe2kZeyH"
	obtcxAddr         = "3e6d9288d04d49585699659aadf3b0a508c47608"
	utxoKey           = "c330431496364497d7257839737b5e4596f5ac06" //"87a9652e9b396545598c0fc72cb5a98848bf93d3"
	toEthAddr         = "0x5cD3143f91a13Fe971043E1e4605C1c23b46bF44"
	ebtcxAddr         = "0x9702640a6b971CA18EFC20AD73CA4e8bA390C910"

	sigs = []string{
		"3044022010efd2114373c5961902333ff29ca79ab033add8bcd1ac90b96260cc37d6556f022041e9f882204259d86d10ff06736dc090ad3d4b1dd555c587c70809ac7e412c0a01",
		"3045022100a9076013a5f0bac46435b7cc5788d23d3054a8d04998ba20228db40174661e68022061ba9c827c08748b5cbb6f744becfbd6344582fb8e34a15255cfe4e9d5fe70c101",
		"30440220737beb8f70082739927a7008282160abdcbd15ca2c671975e8e430f6100645bb0220178f44de5190f83a2fac4871b632a39b3c62f471cd328c8897ce76c95aa52b9a01",
		"3044022007151b9b211ec9ab6ba92bc7f1d92ef8be65024f38c82fe973a2d1d14d8dff8d022006721914f148426ae92baafd5478e3b12cd05a81b99b4f77b9dfd9b4f5d5e5fc01",
		"30440220162dd3953f1cf211eb15d0ca692b36b116abdd8e3608f7158fec44db96cc4a720220428dd1e1f5c02f3de0243d20d6969aa45abc27f93e29b6b39e25e23d3d65002001",
	}

	signers = []string{
		"mxtJn3aRsKLrWRLLAhu2nCBuK2brfazedj",
		"msu1qgtn4FsQh7xDP15ggStwW4yHquUTYE",
		"mzr5T4PmqzNtusmM2S889LUySC9BBbwRzd",
		"mmDSSyis1sjysaCKZ9eK1ww92Vr9S1CNEX",
		"n1WmUbJ4dQfvzvtRaNcjmxH6nADhhv7W8c",
	}

	getNativeFunc = func(args []byte, db *storage.CacheDB) *native.NativeService {
		if db == nil {
			store, _ := leveldbstore.NewMemLevelDBStore()
			db = storage.NewCacheDB(overlaydb.NewOverlayDB(store))
		}

		return native.NewNativeService(db, nil, 0, 0, common.Uint256{0}, 0, args, false)
	}

	setSideChain = func(ns *native.NativeService) {
		side := &side_chain_manager.SideChain{
			Name:         "btc",
			ChainId:      1,
			BlocksToWait: 1,
			Router:       0,
		}
		sink := common.NewZeroCopySink(nil)
		_ = side.Serialization(sink)

		ns.GetCacheDB().Put(utils.ConcatKey(utils.SideChainManagerContractAddress,
			[]byte(side_chain_manager.SIDE_CHAIN), utils.GetUint64Bytes(1)), states.GenRawStorageItem(sink.Bytes()))
	}

	getSigs = func() [][]byte {
		res := make([][]byte, len(sigs))
		for i, sig := range sigs {
			raw, _ := hex.DecodeString(sig)
			res[i] = raw
		}
		return res
	}

	registerRC = func(db *storage.CacheDB) *storage.CacheDB {
		ca, _ := hex.DecodeString(strings.Replace(ebtcxAddr, "0x", "", 1))
		cb := &side_chain_manager.ContractBinded{
			Contract: ca,
			Ver: 0,
		}
		sink := common.NewZeroCopySink(nil)
		cb.Serialization(sink)
		rk, _ := hex.DecodeString(utxoKey)
		db.Put(utils.ConcatKey(utils.SideChainManagerContractAddress, []byte(side_chain_manager.REDEEM_BIND),
			utils.GetUint64Bytes(1), utils.GetUint64Bytes(2), rk), states.GenRawStorageItem(sink.Bytes()))

		return db
	}
)

func TestBTCHandler_MakeDepositProposal(t *testing.T) {
	gh := netParam.GenesisBlock.Header
	mr, _ := chainhash.NewHashFromStr("f68a7646782d6d52f63f94633f4b4fb6cea67d1f14c80fc1fa8f999014c99359")
	gh.MerkleRoot = *mr

	db, err := syncGenesisHeader(&gh)
	if err != nil {
		t.Fatal(err)
	}
	db = registerRC(db)

	txid, _ := chainhash.NewHashFromStr("470284c7e435b28b379901668e4129c408bf1253874317842e2986282265b423")
	rawTx, _ := hex.DecodeString("01000000017811f6394ef9b03e3953ac48de13c69e2d9865fe3f6e17dd42ccfd0c4694b015000000006a47304402206b048cdd7f19be0d3f46a7e785339cf81cf35b33be75297254f41bac7548985202202520d33eaf7a728ed4fab42056e241de5d375378145251e429ade0327c92aa07012102141d092eca49eac51de2760d28cbced212b60efc23fdcbb57304823bb17aa64effffffff0300e1f50500000000220020216a09cb8ee51da1a91ea8942552d7936c886a10b507299003661816c0e9f18b0000000000000000286a266602000000000000000000000000000000145cd3143f91a13fe971043e1e4605c1c23b46bf44180d1024010000001976a9145f35a2cc0318fbc17c4c479964734e7a9f8819d788ac00000000")
	scAddr, _ := hex.DecodeString(strings.Replace(ebtcxAddr, "0x", "", 1))
	handler := NewBTCHandler()

	// wrong proof
	wrongProof := "0100003037db655b09de3449fe60bc0838ef3541e28d3ae31a05093f1bb63e4845a6b102695fd2a687fc1fc368f13227c2bb1b6b0fcac9760936d869a9ba01f8a75f825c5105245effff7f20050000000200000002e0c8d9fb711dd377d0ba8d1c16c154b903432aa8923c87f3f0fd6045be7b8c8a51cf2962a492309e6bd7aa56848b2817c1a760f9fb0823762200d2b286f988b90105"
	proof, _ := hex.DecodeString(wrongProof)
	params := new(ccmcom.EntranceParam)
	params.Height = 0
	params.SourceChainID = 1
	params.Proof = proof
	params.Extra = rawTx
	params.RelayerAddress = acct.Address[:]

	sink := common.NewZeroCopySink(nil)
	params.Serialization(sink)
	ns := getNativeFunc(sink.Bytes(), db)
	setSideChain(ns)
	p, err := handler.MakeDepositProposal(ns)
	assert.Error(t, err)

	// normal case
	proof, _ = hex.DecodeString("0000002012f9b192ec5ec403d9a438dcfdf00f2025118cffafcdc0b26c88218a0d26b6355993c91490998ffac10fc8141f7da6ceb64f4b3f63943ff6526d2d7846768af6ca1b5f5effff7f20000000000200000002db0fcc77abf1be640e42d997a4eeacc4bf7075e791de8d3be21dcb29c337efe323b465222886292e841743875312bf08c429418e660199378bb235e4c78402470105")
	params.Proof = proof

	sink.Reset()
	params.Serialization(sink)
	ns = getNativeFunc(sink.Bytes(), db)

	p, err = handler.MakeDepositProposal(ns)
	assert.NoError(t, err)

	ethAddr, _ := hex.DecodeString(strings.Replace(toEthAddr, "0x", "", 1))
	sink.Reset()
	sink.WriteVarBytes(ethAddr[:])
	sink.WriteUint64(uint64(100000000))

	assert.Equal(t, txid[:], p.CrossChainID)
	assert.Equal(t, txid[:], p.TxHash)
	assert.Equal(t, "unlock", p.Method)
	assert.Equal(t, utxoKey, hex.EncodeToString(p.FromContractAddress))
	assert.Equal(t, uint64(2), p.ToChainID)
	assert.Equal(t, sink.Bytes(), p.Args)
	assert.Equal(t, scAddr[:], p.ToContractAddress)

	utxos, err := getUtxos(ns, 1, utxoKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(utxos.Utxos))
	assert.Equal(t, uint64(100000000), utxos.Utxos[0].Value)
	assert.Equal(t, "470284c7e435b28b379901668e4129c408bf1253874317842e2986282265b423:0", utxos.Utxos[0].Op.String())

	// repeated commit
	sink.Reset()
	params.Serialization(sink)
	ns = getNativeFunc(sink.Bytes(), db)
	p, err = handler.MakeDepositProposal(ns)
	assert.Error(t, err)
}

func TestBTCHandler_MakeTransaction(t *testing.T) {
	rawTx, _ := hex.DecodeString(fromBtcRawTx)
	mtx := wire.NewMsgTx(wire.TxVersion)
	_ = mtx.BtcDecode(bytes.NewBuffer(rawTx), wire.ProtocolVersion, wire.LatestEncoding)

	ns := getNativeFunc(nil, nil)
	_ = addUtxos(ns, 1, 0, mtx)
	registerRC(ns.GetCacheDB())
	setSideChain(ns)

	rk, _ := hex.DecodeString(utxoKey)
	scAddr, _ := hex.DecodeString(strings.Replace(ebtcxAddr, "0x", "", 1))
	r, _ := hex.DecodeString(rdm)

	sink := common.NewZeroCopySink(nil)
	sink.WriteVarBytes([]byte("mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57"))
	sink.WriteUint64(6000)
	sink.WriteVarBytes(r)
	p := &ccmcom.MakeTxParam{
		ToChainID:           1,
		TxHash:              []byte{1},
		Method:              "unlock",
		ToContractAddress:   rk,
		CrossChainID:        []byte{1},
		FromContractAddress: scAddr,
		Args:                sink.Bytes(),
	}

	handler := NewBTCHandler()
	err := handler.MakeTransaction(ns, p, 2)
	assert.NoError(t, err)
	s := ns.GetNotify()[0].States.([]interface{})
	assert.Equal(t, utxoKey, s[1].(string))
}

func TestBTCHandler_MultiSign(t *testing.T) {
	rawTx, _ := hex.DecodeString(fromBtcRawTx)
	mtx := wire.NewMsgTx(wire.TxVersion)
	_ = mtx.BtcDecode(bytes.NewBuffer(rawTx), wire.ProtocolVersion, wire.LatestEncoding)

	ns := getNativeFunc(nil, nil)
	_ = addUtxos(ns, 1, 0, mtx)

	rb, _ := hex.DecodeString(rdm)
	err := makeBtcTx(ns, 1, map[string]int64{"mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57": 6000}, []byte{123},
		2, rb, hex.EncodeToString(btcutil.Hash160(rb)))
	assert.NoError(t, err)
	stateArr := ns.GetNotify()[0].States.([]interface{})
	assert.Equal(t, "makeBtcTx", stateArr[0].(string))
	assert.Equal(t, utxoKey, stateArr[1].(string))

	stxos, err := getStxos(ns, 1, utxoKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stxos.Utxos))
	assert.Equal(t, uint64(10000), stxos.Utxos[0].Value)
	assert.Equal(t, fromBtcTxid+":0", stxos.Utxos[0].Op.String())

	rawTx, _ = hex.DecodeString(stateArr[2].(string))
	_ = mtx.BtcDecode(bytes.NewBuffer(rawTx), wire.ProtocolVersion, wire.LatestEncoding)
	assert.Equal(t, int64(4000), mtx.TxOut[1].Value)

	handler := NewBTCHandler()
	sigArr := getSigs()
	txid := mtx.TxHash()
	// commit no.1 to 4 sig
	for i, sig := range sigArr[:4] {
		msp := ccmcom.MultiSignParam{
			ChainID:   1,
			TxHash:    txid.CloneBytes(),
			Address:   signers[i],
			RedeemKey: utxoKey,
			Signs:     [][]byte{sig},
		}
		sink := common.NewZeroCopySink(nil)
		msp.Serialization(sink)
		ns = getNativeFunc(sink.Bytes(), ns.GetCacheDB())

		err = handler.MultiSign(ns)
		assert.NoError(t, err)
	}

	// repeated submit sig4
	msp := ccmcom.MultiSignParam{
		ChainID:   1,
		TxHash:    txid.CloneBytes(),
		Address:   signers[3],
		RedeemKey: utxoKey,
		Signs:     [][]byte{sigArr[3]},
	}
	sink := common.NewZeroCopySink(nil)
	msp.Serialization(sink)
	ns = getNativeFunc(sink.Bytes(), ns.GetCacheDB())
	err = handler.MultiSign(ns)
	assert.Error(t, err)

	// right sig but wrong address
	msp = ccmcom.MultiSignParam{
		ChainID:   1,
		TxHash:    txid.CloneBytes(),
		Address:   signers[3],
		RedeemKey: utxoKey,
		Signs:     [][]byte{sigArr[4]},
	}
	sink.Reset()
	msp.Serialization(sink)
	ns = getNativeFunc(sink.Bytes(), ns.GetCacheDB())
	err = handler.MultiSign(ns)
	assert.Error(t, err)

	// commit the last right sig
	msp = ccmcom.MultiSignParam{
		ChainID:   1,
		TxHash:    txid.CloneBytes(),
		Address:   signers[4],
		RedeemKey: utxoKey,
		Signs:     [][]byte{sigArr[4]},
	}
	sink.Reset()
	msp.Serialization(sink)
	ns = getNativeFunc(sink.Bytes(), ns.GetCacheDB())
	err = handler.MultiSign(ns)
	assert.NoError(t, err)
	stateArr = ns.GetNotify()[0].States.([]interface{})
	assert.Equal(t, "btcTxToRelay", stateArr[0].(string))
	assert.Equal(t, hex.EncodeToString([]byte{123}), stateArr[4].(string))

	rawTx, err = hex.DecodeString(stateArr[3].(string))
	assert.NoError(t, err)
	err = mtx.BtcDecode(bytes.NewBuffer(rawTx), wire.ProtocolVersion, wire.LatestEncoding)
	assert.NoError(t, err)

	txid = mtx.TxHash()
	utxos, err = getUtxos(ns, 1, utxoKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(utxos.Utxos))
	assert.Equal(t, uint64(4000), utxos.Utxos[0].Value)
	assert.Equal(t, txid.String()+":1", utxos.Utxos[0].Op.String())
}

func syncGenesisHeader(genesisHeader *wire.BlockHeader) (*storage.CacheDB, error) {
	var buf bytes.Buffer
	_ = genesisHeader.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	btcHander := btc.NewBTCHandler()

	sink := new(common.ZeroCopySink)

	h := make([]byte, 4)
	binary.BigEndian.PutUint32(h, 0)
	params := &hscom.SyncGenesisHeaderParam{
		ChainID:       1,
		GenesisHeader: append(buf.Bytes(), h...),
	}
	sink = new(common.ZeroCopySink)
	params.Serialization(sink)

	ns := getNativeFunc(sink.Bytes(), nil)
	err := btcHander.SyncGenesisHeader(ns)
	if err != nil {
		return nil, err
	}

	return ns.GetCacheDB(), nil
}
