# go-blockchain
A simple blockchain made with GO lang


This project was made based on the blockchain example made with Pyhon of this [blog post](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46)

# Endpoints
* **/transactions/new** to create a new transaction to a block
* **/mine** to tell our server to mine a new block.
* **/chain**  to return the full Blockchain
* **/nodes/register** to accept a list of new nodes in the form of URLs.
* **/nodes/resolve** to implement our Consensus Algorithm, which resolves any conflicts—to ensure a node has the correct chain.
