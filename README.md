Requirements for run:  
   MongoDB 
ENV VARS:
   COUNTER_MONGODB_HOST  //Expexted that counter database will be one for all backends, so extra var where hosted counter mongodb instance
   EVENT_MONGODB_HOST   //mongodb  for events srore
   PORT   //port for app

For build  binary  expected that installed golang


Docker image based on scratch so binary should be maked before docker build via
./build_binary.sh



for local usage:
set your host ip var 
example

export HOST_IP=172.17.0.1


then 

./run_local.sh 

This will up mongodb container then build binary then up application container