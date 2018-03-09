#!/bin/bash
# One of the kubernetes cluster's IP
external_url=localhost

# Deploy the echoserver service
kubectl create -f echoserver.yaml

# change the kong service name if necessary
kong_admin_svc=apihub-kong-kong-admin
kong_proxy_svc=apihub-kong-kong-proxy

# Get Kong admin and proxy's port
kong_admin_port=$(kubectl get service $kong_admin_svc -o jsonpath={.spec.ports[0].nodePort})
kong_proxy_port=$(kubectl get service $kong_proxy_svc -o jsonpath={.spec.ports[0].nodePort})

# Remove existing api path
curl -k -X DELETE https://${external_url}:${kong_admin_port}/apis/echoserver

# Register the API
curl -k -i -S -X PUT --url https://${external_url}:${kong_admin_port}/apis \
--data 'name=echoserver' --data 'upstream_url=http://http-echoserver:80' \
--data 'uris=/echoserver' 

sleep 5

# Test the API
curl -k https://${external_url}:${kong_proxy_port}/echoserver
