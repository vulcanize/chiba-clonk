# Vulcanize dxns

Basic node status:

```graphql
{
    getStatus {
        version
        node {
            id
            network
            moniker
        }
        sync {
            latest_block_height
            catching_up
        }
        num_peers
        peers {
            is_outbound
            remote_ip
        }
        disk_usage
    }
}
```

Full node status:

```graphql
{
    getStatus {
        version
        node {
            id
            network
            moniker
        }
        sync {
            latest_block_hash
            latest_block_time
            latest_block_height
            catching_up
        }
        validator {
            address
            voting_power
            proposer_priority
        }
        validators {
            address
            voting_power
            proposer_priority
        }
        num_peers
        peers {
            node {
                id
                network
                moniker
            }
            is_outbound
            remote_ip
        }
        disk_usage
    }
}
```

Get account details:

```graphql
{
    getAccounts(addresses: ["cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094"]) {
        address
        pubKey
        number
        sequence
    balance {
      type
      quantity
    }
  }
}
```


Query bonds:

```graphql
{
  queryBonds(
    attributes: [
      {
        key: "owner"
        value: { string: "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094" }
      }
    ]
  ) {
    id
    owner
    balance {
      type
      quantity
    }
  }
}
```


Get bonds by IDs.

```graphql
{
    getBondsByIds(ids :
    [
        "1c2b677cb2a27c88cc6bf8acf675c94b69051125b40c4fd073153b10f046dd87",
        "c3f7a78c5042d2003880962ba31ff3b01fcf5942960e0bc3ca331f816346a440"
    ])
    {
        id
        owner
        balance{
            type
            quantity
        }
    }
}
```

Query Bonds by Owner 
```graphql 
{
  queryBondsByOwner(ownerAddresses: ["ethm1mfdjngh5jvjs9lqtt9a7y2hlgw8v3syh3hsqzk"])
  {
    owner
    bonds{
      id
      owner
      balance
    	{
        type
        quantity
      }
    }
  }
}
  
```