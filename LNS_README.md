# The Laconic Name Service (LNS)

Laconic Name Service (LNS) is an Ethermint-based blockchains that serves the purpose of mapping human readable names to resources in the Laconic Network, including:

* Watchers
* ...

To ensure the decentralization and permissionless autonomy of the LNS, it is implemented using a bonded auction system that allows the creation of new Authorities (top level namespaces) with resonable protection against squatters and spammers.

## Modules

The following Ethermint modules have been added to this distribution:

* x/auction: governs the sale of top level namespaces
* x/bond: allows the deposit of LNT to the XXX contract which is bonded (during the auction | for the duration of the use of the namespace | forever and ever)
* x/evm: ?
* x/feemarket ?
* x/nameservice: Naming resolver. Used to search for watchers...

## Bonding and Auctions

A person wishing to register a top-level namespace under which to register Laconic Network Watchers must first put LNT into escrow. 

How long are auctions?
What tokens are accepted?
What type of auction and how is bidding resolved?