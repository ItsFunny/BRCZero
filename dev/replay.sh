
start_height=$1
end_height=$2


run() {
  request=$(curl -s --location "http://127.0.0.1:80/api/v1/brc0/rpc_request/$1" | jq -r '.data')
  btchash=$(echo $request | jq -r '.block_hash')
  if [ "$btchash" == "" ]; then
    echo "height: $1 is not exist"
    exit 1
  else
    echo "request$1: $btchash"
    curl --location 'http://localhost:26657' \
    --header 'Content-Type: application/json' \
    --data "{\"jsonrpc\":\"2.0\",\"id\":3,\"method\":\"broadcast_brczero_txs_async\",\"params\":$request}"
  fi
}

for i in `seq $start_height $end_height`
do
   run $i
done
