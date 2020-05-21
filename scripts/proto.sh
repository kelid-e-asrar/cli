#! /bin/bash
if [ -d "contract" ]; then
    cd contract
    git pull origin master
else
    git clone https://github.com/kelid-e-asrar/contract.git contract
    cd contract
fi

for f in $(ls *.proto); do
    echo "$f";
    protoc --go_out=plugins=grpc:. "$f" 2> /dev/null
done


