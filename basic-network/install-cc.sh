#!/bin/bash
set -ev

CCNAME="learning-cc-ch21";

# install chaincode for channelsales1
docker exec cli1 peer chaincode install -n $CCNAME -v 1.0 -p chaincode/go
sleep 1
# instantiate chaincode for channelsales1
docker exec cli1 peer chaincode instantiate -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -v 1.0 -c '{"Args":[""]}' -P "OR ('Sales1Org.member','CustomerOrg.member')"
sleep 10

# enrollUser | INPUT: userID, user_key | OUTPUT: 200 or 400
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"enrollUser","Args":["user1", "user1_key"]}'
sleep 1
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"enrollUser","Args":["user2", "user2_key"]}'
sleep 1
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"enrollUser","Args":["user3", "user3_key"]}'
sleep 1

# updateUserToken | INPUT: userID, Token | OUTPUT: 200 or 400
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"updateUserToken","Args":["user1", "10"]}'
sleep 1

# getUserInfo | INPUT: userID | OUTPUT: userID, user_key, Token
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user1"]}'
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user2"]}'
sleep 1

# updateGlobalAccuracy | INPUT: round | OUTPUT: 200 or 400
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"updateGlobalAccuracy","Args":["0", "99.9"]}'
sleep 1

# getGlobalAccuracy | INPUT: round | OUTPUT: globalAccuracy
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getGlobalAccuracy","Args":["0"]}'
sleep 1

# quitUser | INPUT: userID | OUTPUT: 200 or 400
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"quitUser","Args":["user2"]}'
sleep 1

# 사용자가 삭제되었는지 확인하는 용도 오류 발생 한다.
# docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserInfo","Args":["user2"]}'
# sleep 1


# uploadWeight | INPUT: userID, weight, round | OUTPUT: 200 or 400
# json input format
# "{\"conv0weight\":\"[0.1, 0.2, 0.3]\",\"conv0bias\":\"[0.1, 0.2, 0.3]\",\"conv3weight\":\"[0.1, 0.2, 0.3]\",\"conv3bias\":\"[0.1, 0.2, 0.3]\",\"conv6weight\":\"[0.1, 0.2, 0.3]\",\"conv6bias\":\"[0.1, 0.2, 0.3]\",\"fc0weight\":\"[0.1, 0.2, 0.3]\",\"fc0bias\":\"[0.1, 0.2, 0.3]\",\"fc2weight\":\"[0.1, 0.2, 0.3]\",\"fc2bias\":\"[0.1, 0.2, 0.3]\"}"
docker exec cli1 peer chaincode invoke -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"uploadWeight","Args":["user1","{\"conv0weight\":\"[0.1, 0.2, 0.3]\",\"conv0bias\":\"[0.1, 0.2, 0.3]\",\"conv3weight\":\"[0.1, 0.2, 0.3]\",\"conv3bias\":\"[0.1, 0.2, 0.3]\",\"conv6weight\":\"[0.1, 0.2, 0.3]\",\"conv6bias\":\"[0.1, 0.2, 0.3]\",\"fc0weight\":\"[0.1, 0.2, 0.3]\",\"fc0bias\":\"[0.1, 0.2, 0.3]\",\"fc2weight\":\"[0.1, 0.2, 0.3]\",\"fc2bias\":\"[0.1, 0.2, 0.3]\"}","0"]}'
sleep 3

# getUserWeight | INPUT: userID, round | OUTPUT: 해당 round에 userID가 등록한 가중치 반환
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getUserWeight","Args":["user1","0"]}'
sleep 3

# aggregation | INPUT: round | OUTPUT: 해당 round의 모든 모델 가중치를 aggregation하여 반환한다.
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"aggregation","Args":["0"]}'

# getAlluserID | INPUT: | OUPUT: 블록체인에 참여한 모든 사용자 ID
docker exec cli1 peer chaincode query -o orderer1.chainbee.com:7050 -C channelsales1 -n $CCNAME -c '{"function":"getAllUserID","Args":[]}'



