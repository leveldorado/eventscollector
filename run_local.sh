#!/bin/bash

./install_mongodb_container.sh
./build_binary.sh
sudo docker build -t eventscollector  .
sudo docker run --name eventscollector  -p 13000:13000 -e 'COUNTER_MONGODB_HOST=http:localhost:27017'  -e 'EVENT_MONGODB_HOST=http:localhost:27017'  eventscollector 