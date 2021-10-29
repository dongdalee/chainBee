#!/bin/bash
set -ev

CCNAME="learning-cc-ch";

# install chaincode for channelsales1
docker exec cli1 peer chaincode install -n $CCNAME -v 1.0 -p chaincode/go
sleep 1
# instantiate chaincode for channelsales1
docker exec cli1 peer chaincode instantiate -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -v 1.0 -c '{"Args":[""]}' -P "OR ('Sales1Org.member','CustomerOrg.member')"
sleep 10

# invoke chaincode for channelsales1
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"enrollUser","Args":["user1", "user1_key"]}'
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"enrollUser","Args":["user2", "user2_key"]}'
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"updateUserToken","Args":["user1", "10"]}'
sleep 3
# query chaincode for channelsales1
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user1"]}'
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user2"]}'
sleep 1

# invoke chaincode for channelsales1
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"updateGlobalAccuracy","Args":["0", "99.9"]}'
sleep 3
# query chaincode for channelsales1
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getGlobalAccuracy","Args":["0"]}'
sleep 1
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"quitUser","Args":["user2"]}'
sleep 3
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"uploadWeight","Args":["user1",`{
    "conv0weight":"[0.1, 0.2, 0.3]",
    "conv0bias":"[0.2, 0.2, 0.3]",
    "conv3weight":"[0.3, 0.2, 0.3]",
    "conv3bias":"[0.4, 0.2, 0.3]",
    "conv6weight":"[0.5, 0.2, 0.3]",
    "conv6bias":"[0.6, 0.2, 0.3]",
    "fc0weight":"[0.7, 0.2, 0.3]",
    "fc0bias":"[0.8, 0.2, 0.3]",
    "fc2weight":"[0.9, 0.2, 0.3]",
    "fc2bias":"[1.0, 0.2, 0.3]"
}` ,"0"]}'

docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"aggregation","Args":["0"]}'

# docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user2"]}'