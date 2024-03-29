# WEB3-CLI

[![test](https://github.com/islishude/web3-cli/actions/workflows/test.yaml/badge.svg)](https://github.com/islishude/web3-cli/actions/workflows/test.yaml)

## Install

You can download the binary in the latest [release page](https://github.com/islishude/web3-cli/releases/latest).

or you have docker

```
docker pull ghcr.io/islishude/web3-cli
```

or if you have golang installed

```sh
# install with the laest tag
go install github.com/islishude/web3-cli@latest
# install with the latest commit
go install github.com/islishude/web3-cli@main
```

## Usage

**Do a simple jsonrpc call**

the default rpc endpoint is your local: `http://localhost:8545`

your command is like `web3-cli jsonrpc_method [jsonrpc_param...]`

```console
$ web3-cli web3_clientVersion
"Geth/v1.12.1-unstable-60070fe5-20230805/darwin-amd64/go1.20.7"
$ web3-cli eth_getBlockByNumber 1 false
{
    "baseFeePerGas": "0x342770c0",
    "difficulty": "0x0",
    "extraData": "0xd983010c01846765746888676f312e32302e378664617277696e",
    "gasLimit": "0xafa5bd",
    "gasUsed": "0x5208",
    "hash": "0xf5c52795b6fb4b69601e8613c436e07d964b33fb97b5b2b8670a106e0617bf6e",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "0x0000000000000000000000000000000000000000",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x1",
    "parentHash": "0xd7bcb282488afadddab1b33e1d5906d1cf367bfc6249b6eaf140961e986621c2",
    "receiptsRoot": "0xf78dfb743fbd92ade140711c8bbc542b5e307f0ab7984eff35d751969fe57efa",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x2ad",
    "stateRoot": "0x05b0dc6527b6d706aa6400866014f8f0bc4996ac8f7b9c4d7b4801fccfad139f",
    "timestamp": "0x64fc7d7a",
    "totalDifficulty": "0x0",
    "transactions": [
        "0xfae530d45ec70be05c68aaa0628588f0055f4aae360a8c090c2591e8691a47ca"
    ],
    "transactionsRoot": "0x8b2e1c5c5426fb942d176fac9186a05911c41867709a9fa43b05e2051a64f712",
    "uncles": [],
    "withdrawals": [],
    "withdrawalsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
}
```

web3-cli supports built-in chain config as well.

```console
$ web3-cli --chain eth eth_blockNumber
"0x1142cb9"
$ web3-cli --chain bsc eth_chainId
"0x38"
```

You can run `web3-cli chains` to get all built-in chain configs.

Of course, you can also use your own rpc endpoint

```console
$ web3-cli --rpc https://rpc.ankr.com/arbitrum eth_chainId
"0xa4b1"
```

**Call a contract with human-readable params**

```console
$ export USDT_TOKEN=0xdac17f958d2ee523a2206206994597c13d831ec7
$ web3-cli --chain eth --abi-name erc20 --call-to $USDT_TOKEN symbol
[
    "USDT"
]
$ web3-cli --chain eth --call-to $USDT_TOKEN balanceOf 0x0000000000000000000000000000000000000000
[
    "0xbf17d330"
]
```

the `--abi-name` could be a built-in abi name, you can `web3-cli abis` to get the built-in abi list.

and it also could be a url and file path.

```
$ web3-cli --chain eth --abi-path https://http-server/abi/abi.json
$ web3-cli --chain eth --abi-path local/path/to/abi.json
```

if you don't provide `--abi-name` or `--abi-path`, web3-cli can fetch the abi from explorer api automatically

```console
$ web3-cli --chain eth --call-to $USDT_TOKEN getOwner
[
    "0xc6cde7c39eb2f0f0095f41570af89efc2c1ea828"
]
```

you may need to provide your explorer api endpoint as well if you provide a custom rpc

```console
$ # not required if use a built-in chain
$ web3-cli --chain eth --rpc https://my-own-rpc.com
$ web3-cli --rpc https://my-own-rpc.com --explorer-api https://my-custom-explorer.com/api
```

What about a complex abi type parameter, like array and tuple?

**bytes**

```solidity
function logbyt(bytes memory) public returns (bytes memory)
```

you must use a hex string

```
0x776562332d636c6920697320736f20636f6f6c21
```

**array and slice**

```solidity
function add(uint256[] calldata items) public pure returns (uint256)
```

json array is valid

```json
["0x1", 100]
```

**tuple(struct)**

```solidity
struct Payment {
    address payable to;
    uint256 value;
}

function transfer(Payment calldata item) external payable returns (bool success) {
    return item.to.send(msg.value);
}
```

json array is valid (you can use it in Remix)

```json
["0x0000000000000000000000000000000000000000", "0x1"]
```

json object is valid as well!

the key is the field name, and it is not case sensitive.

```json
{ "to": "0x0000000000000000000000000000000000000000", "value": "0x1" }
{ "To": "0x0000000000000000000000000000000000000000", "Value": "0x1" }
```

## Built-in tools

**Create a new address with optional prefix and suffix**

```console
$ web3-cli tools new-random-address --prefix 00 --suffix 00
{
    "Address": "0x00BbEA9643992F70879a1f4009a18c82B14f1A00",
    "PrivateKey": "0xc44591dd44df800b1f9$REDACTED$be239403741b2196ca28",
    "PublicKey": "0x044708bac89fd0c4f595b0808c068350cf889d7aeea67848498997458c0541da955a7b6c8aafd5e8fbe00edfabb2d47da9fee680b711b2ce90401145ffce90defe"
}
```

**Decode raw transactions**

NOTE: OP system transactions are not supported.

```console
$ web3-cli tools decode-raw-tx 0x...REDACTED...
{
    "type": "0x2",
    "chainId": "0x1",
    "nonce": "0x1018",
    "to": "0x388c818ca8b9251b393131c08a736a67ccb19297",
    "gas": "0x565f",
    "gasPrice": null,
    "maxPriorityFeePerGas": "0x0",
    "maxFeePerGas": "0x50c81c068",
    "value": "0x475c0000a4d611",
    "input": "0x",
    "accessList": [],
    "v": "0x1",
    "r": "0xa5a963fee24751d6f54656be527c699584ef47aade8a677c806f119b02b0daf7",
    "s": "0x257af2dbf67d167dfbf4f4acf032b43ca4fa8b8f06ae8b09f16bd2cc57817ce4",
    "yParity": "0x1",
    "hash": "0x048117077d33c6f3670d601de5525dda41374719162cb0ed252726559d4ffe70"
}
```

## Contribution

**How to add your chain to built-in list?**

Add it to [internal/chains/list.go](./internal/chains/list.go) file.

for example:

```go
var example = &Chain{
	Name:     "example",
	Id:       111000111,
	Endpoint: "https://jsonrpc-endpoint.example",
	Explorer: "https://full-explorer-api-endpoint.example/api",
}
```

then append it to `Builtin` slice.

you also need to add a test case for it in [internal/chains/explorer_api_test.go](./internal/chains/explorer_api_test.go)

**How to add an ABI to built-in list?**

Add a valid json file to `internal/abis/abi` directory.
