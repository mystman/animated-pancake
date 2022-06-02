#!/bin/sh

# minikube tunnel requires
echo "Some tests with Curl:"
echo " (i) Please note that 'minikube tunnel' must be running to support LoadBalancer on localhost\n"
read -t 2 -p ""


echo "POST to /network"
echo "======================================================================"
curl -X POST -H "Content-Type: application/json" -d '{ "subnetName": "testSubnet", "ipRange": "10.2.1.0/24" }' http://localhost:6543/v1/network


read -t 3 -p ""
echo "\n\n"
echo "GET with ID to /"
echo "======================================================================"
curl -X GET http://localhost:6543/v1/1


read -t 3 -p ""
echo "\n\n"
echo "PUT to /"
echo "======================================================================"
curl -X PUT -H "Content-Type: application/json" -d '{ "subnetName": "CHANGEDSubnet", "ipRange": "10.9.7.7/24" }' http://localhost:6543/v1/2


read -t 3 -p ""
echo "\n\n"
echo "GET (all) to /"
echo "======================================================================"
curl -X GET http://localhost:6543/v1/


read -t 3 -p ""
echo "\n\n"
echo "DELETE to /"
echo "======================================================================"
curl -X DELETE http://localhost:6543/v1/1
