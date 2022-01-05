This document describes the Laconic Name Service (LNS).

# The Laconic Name Service (LNS)

Laconic Name Service (LNS) is a permissionless blockchain that maps human readable names to content addressable resources in the Laconic Network, including:

* Watchers 
* Validators
* Service providers responding to requesters
* Responder contracts
* Other platform assets

To ensure the decentralization and permissionless autonomy of the LNS, it is implemented using a bonded auction system that allows the creation of new Authorities (top level namespaces) with reasonable protection against squatters and spammers.

## Modules

The following custom modules have been added to this application:

* x/auction: A stand-alone module that provides [sealed-bid, second-price auctions](https://en.wikipedia.org/wiki/Vickrey_auction)).
* x/bond: allows the deposit (bond) of the blockchain's token into an escrow account.
* x/nameservice: Naming resolver. Takes parameterized search criteria and returns the addresses of matched resources. It is used, for example, to resolve human readable names to Watcher addresses.  

The `x/nameservice` module uses `x/auction` to govern the sale of top level namespaces called Authorities. Users pay rent out of their bond to retain the Authority. Storing content addressable data likewise requires rent paid out of the bond.

## Bonding and Auctions

When registering names or storing data, a time-based fee is charged, and that is drawn from your bond (escrow account).

### Create a Bond

1. Create a bond .
2. Register a name.
3. Add metadata (eg. watcher registration, metadata, contract data).

Payments for #2 & #3 are drawn from the bond. If the bond runs out, data is not served, and the name registration (after the grace period has lapsed) is lost.

### How an Auction works

Bidders have accounts on the LNS chain. Bidders bid to reserve an Authority. Auctions occur in two phases: Commit and Reveal.

In the **Commit** phase, bidders send a hash of their bid so as not to reveal the sum of the bid to others who are bidding. The Commit phase is timeboxed. 

In the **Reveal** phase, when bidding is over, bidders send their bids to be proven against the hash that was previously submitted. At the end of the Reveal phase, the blockchain picks the highest bidder as the winner. 

The winning bidder then pays the sum represented by the next highest bid to secure the auction.

Some further details govern auctions:
    - There is a minimum bid.
    - Revealing is optional.
    - Bidders pay a fee to bid. The first component of this fee is returned if the bidder reveals during the Reveal phase.
    - Bidders who don't reveal lose the second fee component.

## How do Bonds work?

An "Authority bond" is what allows you to set and look up names within an Authority (Top level name).

Specifically, these actions are:
- Withdraw bond (account).
- Cancel bond.
- Associate/Disassociate/Reassociate bond with records.

# Installation instructions

## Build the Chain

These instructions have been tested on Ubuntu and Alpine Linux.

The following command builds the Ethermint daemon and places the binary in the `build` directory.

```
make build
```

## Setup the Chain

You only need to follow these steps before running the chain for the first time.

1. Add the root key:
   ```
   ./build/ethermintd keys add root
   ```
   Keep a note of the keyring passphrase if you set it.
2. Init the chain:
   ```
   ./build/ethermintd init test-moniker --chain-id ethermint_9000-1
   ```
3. Add genesis account:
   ```
   ./build/ethermintd add-genesis-account $(./build/ethermintd keys show root -a) 1000000000000000000aphoton,1000000000000000000stake
   ```
4. Make a genesis tx:
   ```
   ./build/ethermintd gentx root 1000000000000000000stake --chain-id ethermint_9000-1 
   ```
5. Collect gentxs:
   ```
   ./build/ethermintd collect-gentxs
   ```

The chain can now be started using:

```
./build/ethermintd start
```

# Usage

You can see some example queries in the form of GraphQL in the [registry client](https://github.com/vulcanize/dxns-registry-client/blob/main/src/registry_client.js). 

You can find test cases for the bond, auction, and nameservice modules in the corresponding [*.test.js](https://github.com/vulcanize/dxns-registry-client/tree/main/src) files for the registry client.
