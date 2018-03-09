#!/bin/bash
# One of the kubernetes cluster's IP
external_url=127.0.0.1

# Deploy a jupyter nb service
kubectl delete -f jupyternb.yaml
kubectl create -f jupyternb.yaml

# change the kong service name if necessary
kong_admin_svc=apihub-kong-kong-admin
kong_proxy_svc=apihub-kong-kong-proxy

# Get Kong admin and proxy's port
kong_admin_port=$(kubectl get service $kong_admin_svc -o jsonpath={.spec.ports[0].nodePort})
kong_proxy_port=$(kubectl get service $kong_proxy_svc -o jsonpath={.spec.ports[0].nodePort})

# Remove existing api path
curl -k -X DELETE https://${external_url}:${kong_admin_port}/apis/notebook

# Register the API
# Note the strip_uri needs to be set to false
curl -k -i -S -X PUT --url https://${external_url}:${kong_admin_port}/apis \
--data 'name=notebook' --data 'upstream_url=http://jupyter-notebook:8888/' \
--data 'uris=/notebook/' --data 'strip_uri=false'

# Test the API
echo "Visiting https://${external_url}:${kong_proxy_port}/notebook"
echo "And use 'mytoken' as the token to login"
