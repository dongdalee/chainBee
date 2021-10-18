# Create the channel1
docker exec cli1 peer channel create -o orderer1.chainbee.com:7050 -c channelsales1 -f /etc/hyperledger/configtx/channel1.tx

# Join peer0.sales1.chainbee.com to the channel and Update the Anchor Peers in Channel1
docker exec cli1 peer channel join -b channelsales1.block
docker exec cli1 peer channel update -o orderer1.chainbee.com:7050 -c channelsales1 -f /etc/hyperledger/configtx/Sales1Organchors.tx

# Join peer1.sales1.chainbee.com to the channel
docker exec -e "CORE_PEER_ADDRESS=peer1.sales1.chainbee.com:7051" cli1 peer channel join -b channelsales1.block

# Join peer0.customer.chainbee.com to the channel and update the Anchor Peers in Channel1
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer0.customer.chainbee.com:7051" cli1 peer channel join -b channelsales1.block
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer0.customer.chainbee.com:7051" cli1 peer channel update -o orderer1.chainbee.com:7050 -c channelsales1 -f /etc/hyperledger/configtx/CustomerOrganchorsChannel1.tx

# Join peer0.customer.chainbee.com to the channel
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer1.customer.chainbee.com:7051" cli1 peer channel join -b channelsales1.block

# Create the channel2
docker exec cli2 peer channel create -o orderer1.chainbee.com:7050 -c channelsales2 -f /etc/hyperledger/configtx/channel2.tx

# Join peer0.sales2.chainbee.com to the channel and Update the Anchor Peers in Channel1
docker exec cli2 peer channel join -b channelsales2.block
docker exec cli2 peer channel update -o orderer1.chainbee.com:7050 -c channelsales2 -f /etc/hyperledger/configtx/Sales2Organchors.tx

# Join peer1.sales2.chainbee.com to the channel
docker exec -e "CORE_PEER_ADDRESS=peer1.sales2.chainbee.com:7051" cli2 peer channel join -b channelsales2.block

# Join peer0.customer.chainbee.com to the channel and update the Anchor Peers in Channel1
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer0.customer.chainbee.com:7051" cli2 peer channel join -b channelsales2.block
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer0.customer.chainbee.com:7051" cli2 peer channel update -o orderer1.chainbee.com:7050 -c channelsales2 -f /etc/hyperledger/configtx/CustomerOrganchorsChannel2.tx

# Join peer0.customer.chainbee.com to the channel
docker exec -e "CORE_PEER_LOCALMSPID=CustomerOrg" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/customer.chainbee.com/users/Admin@customer.chainbee.com/msp" -e "CORE_PEER_ADDRESS=peer1.customer.chainbee.com:7051" cli2 peer channel join -b channelsales2.block
