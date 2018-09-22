# PiButterfly: Provably Secure Onion Routing

## Research
Abstract: There have been many onion routing protocols that have either been provably secure or very efficient. We have protocol PiButterfly provided by faculty mentors, which is provably secure. I have programmed and implemented this security protocol exactly how it would be run and plotted the data. We simulated the protocol with different parameters so see things such as scalability, practicality, and efficiency. We closely analyzed what parameters you could forgo for higher security(latency, space, redundancy) and what you don’t have to forgo. We compared the status quo to our own protocol and showed that much efficiency wasn't given up for provably secure results. PiButterfly is one of the only security protocols that are statistically private from the active adversary, has low communication complexity, is fault tolerant and shown by implementation/results, is practical.

Link:(will be updated when it is published)

## Implementation

There are many sections. The program is seperated by category with the following components:

* batching: This is where the batching of the onions occurs. Each participant peels their onions here, choses to drop onions if the protocol desires, and organizes them for sending. This part would be parallized in the actual program.

* crypto: This uses the github.com/dedis/kyber package with protobuffer functions. This is the simple package for encrypyting

* generation: This package creates the checkpoint nonces and generates the path for each onion. This is used for creating the onions

* onionDS: The main data structure for the implementation. Creates an onion that can be encrypted and decrypted properly and forms layers and procs layers.

* onionForming: This is where the users would create the onion and how they would create their merging onions/dummy onions. They also use the nonces from generation.

* piButterfly: This is a package dedicated to the specific implmentation of PiButterfly. All of the other packages can be used to represet PiTree, a less efficient provably secure onion routing protocol. This changes the methods slightly to account for PiButterfly.

* supportCode: This package has basic functions that dont actually complete tasks designed for the protocol, but are needed to run certain aspects(e.g. find an index of an element in an array)

## Running the protocol
To run the protocol, simply install Go, and in the directory "github.com/czebos/pi" run 
  
   go run protocol.go
   
and one run of the protocol will start.

You may change the parameters at the top of the main function, but the standard is:

* lambda - 8
* particpant size - 16
* epsilon - .00005
* corruptionRate - 0
* message - "Security is Great!"

## Known Bugs

* Unfortuantely, when the protocol is run, when onions are dropped, certain checkpoints are never reached. This causes an error, because people expect those onions. This will cause the protocol to go over the threshold and force the protocol to abort. This is not important for the actual simulation, since it doesnt affect time. To account for this, I add in a factor that increases the threshold to account for this
