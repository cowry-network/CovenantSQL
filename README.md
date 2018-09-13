<p align="center">
    <img src="logo/covenantsql_horizontal.png"
        height="130">
</p>
<p align="center">
    <a href="https://goreportcard.com/report/github.com/CovenantSQL/CovenantSQL">
        <img src="https://goreportcard.com/badge/github.com/CovenantSQL/CovenantSQL?style=flat-square"
            alt="Go Report Card"></a>
    <a href="https://codecov.io/gh/CovenantSQL/CovenantSQL">
        <img src="https://codecov.io/gh/CovenantSQL/CovenantSQL/branch/develop/graph/badge.svg"
            alt="Coverage"></a>
    <a href="https://travis-ci.org/CovenantSQL/CovenantSQL">
        <img src="https://travis-ci.org/CovenantSQL/CovenantSQL.png?branch=develop"
            alt="Build Status"/></a>
    <a href="https://opensource.org/licenses/Apache-2.0">
        <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"
            alt="License"></a>
    <a href="https://godoc.org/github.com/CovenantSQL/CovenantSQL">
        <img src="https://img.shields.io/badge/godoc-reference-blue.svg"
            alt="GoDoc"></a>
    <a href="https://twitter.com/intent/follow?screen_name=CovenantLabs">
        <img src="https://img.shields.io/twitter/url/https/twitter.com/fold_left.svg?style=social&label=Follow%20%40CovenantLabs"
            alt="follow on Twitter"></a>
</p>

## 

CovenantSQL is a decentralized, crowdsourcing SQL database on blockchain. with Features:

- **SQL**: most SQL-92 support.
- **Decentralize**: decentralize with our consensus algorithm DH-RPC & Kayak.
- **Privacy**: access with granted permission and Encryption Pass.
- **Immutable**: query history in CovenantSQL is immutable and trackable.

We believe [On the next Internet, everyone should have a complete **Data Rights**](https://medium.com/@covenant_labs/covenantsql-the-sql-database-on-blockchain-db027aaf1e0e)

#### One line makes App to ĐApp
```go
sql.Open("CovenantSQL", dbURI)
```


## Arch

![CovenantSQL 3 Layer design](logo/arch.png)

1. **Global Consensus Layer** (the main chain, the middle ring in the architecture diagram):
    - There will only be one main chain throughout the network.
    - Mainly responsible for database Miner and the user’s contract matching, transaction settlement, anti-cheating, shard chain lock hash and other global consensus matters.
1. **SQL Consensus Layer** (shard chain, rings on both sides):
    - Each database will have its own separate shard chain.
    - Mainly responsible for: the signature, delivery and consistency of the various Transactions of the database. The data history of the permanent traceability is mainly implemented here, and the hash lock is performed in the main chain.
1. **Datastore Layer** (database engine with SQL-92 support):
    - Each Database has its own independent distributed engine.
    - Mainly responsible for: database storage & encryption, query processing & signature, efficient indexing.


## Tech

#### Network Stack

  - [DH-RPC](rpc/) = TLS - Cert + DHT.
  <img src="logo/DH-RPC-Layer.png" width=350>

    - [**E**nhanced **TLS**](https://github.com/CovenantSQL/research/wiki/ETLS(Enhanced-Transport-Layer-Security)): the Transport Layer Security.
  
#### Test Tools
  -  [(**G**lobal **N**etwork **T**opology **E**mulator)](https://github.com/CovenantSQL/GNTE) is used for network emulating.


#### Connector

CovenantSQL is Still Under Construction(U know..). Test net will be released till Oct. 

Watch us or [![follow on Twitter](https://img.shields.io/twitter/url/https/twitter.com/fold_left.svg?style=social&label=Follow%20%40CovenantLabs)](https://twitter.com/intent/follow?screen_name=CovenantLabs) we will *Wake U Up When September Ends*

- [Golang](client/)
- [Java](https://github.com/CovenantSQL/covenant-connector)
- Coding for more……

## Contact

- [mail us](mailto:webmaster@covenantsql.io)
- <a href="https://twitter.com/intent/follow?screen_name=CovenantLabs">
          <img src="https://img.shields.io/twitter/url/https/twitter.com/fold_left.svg?style=social&label=Follow%20%40CovenantLabs"
              alt="follow on Twitter"></a>



