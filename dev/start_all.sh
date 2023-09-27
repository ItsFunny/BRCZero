#!/bin/bash

./kill_all.sh

set -e
set -o errexit
set -a
set -m

echo "************* Start brczero node... *************"
./testnet.sh -s -i -n 4

echo "************* Start btc node... *************"
rm -r ./bitcoin-data/regtest
docker-compose -f bitcoin.yml up -d
cp -r ./bitcoin-data/wallets ./bitcoin-data/regtest/
sleep 5
docker exec -it local_bitcoin_node bitcoin-cli loadwallet testwallet_01
docker exec -it local_bitcoin_node bitcoin-cli -rpcwallet=testwallet_01 getwalletinfo
docker exec -it local_bitcoin_node bitcoin-cli generatetoaddress 120 bcrt1qd28jewrz9y9hpl328em5fpljvgarucgcxf7fjt
docker exec -it local_bitcoin_node bitcoin-cli -rpcwallet=testwallet_01 getwalletinfo

echo "************* Start ord... *************"
cd ~/rust/BRC20S
rm -rf ./_cache
nohup ./target/debug/ord \
 --log-level=INFO \
 --data-dir=./_cache \
 --rpc-url=http://localhost:18443 \
 --regtest \
 --bitcoin-rpc-user bitcoinrpc \
 --bitcoin-rpc-pass bitcoinrpc \
 server >/dev/null 2>&1 &

