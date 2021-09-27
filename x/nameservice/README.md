# Build chain
```bash 
# it will create binary in build folder with `ethermintd` 
$ make build 
```
# Setup Chain
```bash
./build/ethermintd keys add root 
./build/ethermintd init test-moniker --chain-id ethermint_9000-1
./build/ethermintd add-genesis-account $(./build/ethermintd keys show root -a) 1000000000000000000aphoton,1000000000000000000stake
./build/ethermintd gentx root 1000000000000000000stake --chain-id ethermint_9000-1 
./build/ethermintd collect-gentxs
./build/ethermintd start
```

## Get Params 
```bash
$ ./build/ethermintd q nameservice params -o json | jq .
{
  "params": {
    "record_rent": {
      "denom": "stake",
      "amount": "1000000"
    },

    "record_rent_duration": "31536000s",
    "authority_rent": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_rent_duration": "31536000s",
    "authority_grace_period": "172800s",
    "authority_auction_enabled": false,
    "authority_auction_commits_duration": "86400s",
    "authority_auction_reveals_duration": "86400s",
    "authority_auction_commit_fee": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_auction_reveal_fee": {
      "denom": "stake",
      "amount": "1000000"
    },
    "authority_auction_minimum_bid": {
      "denom": "stake",
      "amount": "5000000"
    }
  }

```

## Create (Set) Record 
> First you have to Create bond 
```bash
$ ./build/ethermintd tx nameservice set ~/Desktop/examples/records/example1.yml 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d --from root --chain-id ethermint_9000-1 --yes -o json
{
  "height": "0",
  "txhash": "BA44ABE1194724694E7CB290F9F3121DB4E63E1A030D95CB84813EEA132CF95F",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}
```

## Get records list 
```bash
$ ./build/ethermintd q nameservice list -o json | jq . 
{
  "records": [
    {
      "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
      "bond_id": "95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d",
      "create_time": "2021-09-27T07:23:25.558111606Z",
      "expiry_time": "2022-09-27T07:23:25.558111606Z",
      "deleted": false,
      "owners": [],
      "attributes": "eyJhdHRyMSI6InZhbHVlMSIsImF0dHIyIjoidmFsdWUyIiwibGluazEiOnsiLyI6IlFtU251V214cHRKWmRMSnBLUmFyeEJNUzJKdTJvQU5WcmdicjJ4V2JpZTliMkQifSwibGluazIiOnsiLyI6IlFtUDhqVEcxbTlHU0RKTENiZVdoVlNWZ0V6Q1BQd1hSZENSdUp0UTVUejlLYzkifX0="
    }
  ],
  "pagination": null
}

```

## Get record by id 
```bash
$ ./build/ethermintd q nameservice get bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae -o json | jq .
{
  "record": {
    "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
    "bond_id": "95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d",
    "create_time": "2021-09-27T07:23:25.558111606Z",
    "expiry_time": "2022-09-27T07:23:25.558111606Z",
    "deleted": false,
    "owners": [],
    "attributes": "eyJhdHRyMSI6InZhbHVlMSIsImF0dHIyIjoidmFsdWUyIiwibGluazEiOnsiLyI6IlFtU251V214cHRKWmRMSnBLUmFyeEJNUzJKdTJvQU5WcmdicjJ4V2JpZTliMkQifSwibGluazIiOnsiLyI6IlFtUDhqVEcxbTlHU0RKTENiZVdoVlNWZ0V6Q1BQd1hSZENSdUp0UTVUejlLYzkifX0="
  }
}
```

## Reserver authority name 
```bash
 ./build/ethermintd tx nameservice reserve-authority hello --from root --chain-id ethermint_9000-1 --owner $(./build/ethermintd key
s show root -a) -y -o json | jq .
{
  "height": "0",
  "txhash": "7EC19157AC89279DEBE840EA3384FC95D1E2A0931C27746CA42AC23AE285B7ED",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}

```
## Query Whois for name authority 
```bash
 ./build/ethermintd q nameservice whois hello -o json | jq .
{
  "name_authority": {
    "owner_public_key": "Au3hH1tzL1KgZfXfA71jGYSe5RV9Wg95kwhBWs8V+N+h",
    "owner_address": "ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk",
    "height": "174",
    "status": "active",
    "auction_id": "",
    "bond_id": "",
    "expiry_time": "2021-09-29T07:34:36.304545965Z"
  }
}

```
## Query the nameservice module balance 
```bash
$ ./build/ethermintd q nameservice  balance -o json | jq .
{
  "balances": [
    {
      "account_name": "record_rent",
      "balance": [
        {
          "denom": "stake",
          "amount": "1000000"
        }
      ]
    }
  ]
}

```

## add bond to the authority 
```bash
$ ./build/ethermintd tx nameservice authority-bond [Authority Name] [Bond ID ]  --from root --chain-id ethermint_9000-1  -y -o json | jq .  
$ ./build/ethermintd tx nameservice authority-bond hello 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d  --from root --chain-id ethermint_9000-1  -y -o json | jq .  
 ```

## Query the records by associate bond id 
```bash
./build/ethermintd q nameservice query-by-bond 95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d -o json | jq .
{
  "records": [
    {
      "id": "bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae",
      "bond_id": "95f68b1b862bfd1609b0c9aaf7300287b92fec90ac64027092c3e723af36e83d",
      "create_time": "2021-09-27T08:25:32.893155609Z",
      "expiry_time": "2022-09-27T08:25:32.893155609Z",
      "deleted": false,
      "owners": [],
      "attributes": "eyJhdHRyMSI6InZhbHVlMSIsImF0dHIyIjoidmFsdWUyIiwibGluazEiOnsiLyI6IlFtU251V214cHRKWmRMSnBLUmFyeEJNUzJKdTJvQU5WcmdicjJ4V2JpZTliMkQifSwibGluazIiOnsiLyI6IlFtUDhqVEcxbTlHU0RKTENiZVdoVlNWZ0V6Q1BQd1hSZENSdUp0UTVUejlLYzkifX0="
    }
  ],
  "pagination": null
}

```

## Renew a record 
> When a record is expires , needs to renew record 
```bash
./build/ethermintd tx nameservice renew-record bafyreih7un2ntk235wshncebus5emlozdhdixrrv675my5umb6fgdergae --from root --chain-id ethermint_9000-1 

```
## Set the authority name 
```bash
$ ./build/ethermintd tx nameservice set-name wrn://hello/test test_hello_cid  --from root --chain-id ethermint_9000-1 -y -o json | jq .
{
  "height": "0",
  "txhash": "66A63C73B076EEE9A2F7605354448EDEB161F0115D4D03AF68C01BA28DB97486",
  "codespace": "",
  "code": 0,
  "data": "",
  "raw_log": "[]",
  "logs": [],
  "info": "",
  "gas_wanted": "0",
  "gas_used": "0",
  "tx": null,
  "timestamp": ""
}
```