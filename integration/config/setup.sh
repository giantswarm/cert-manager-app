#!/usr/bin/env bash

# label the master so it matches our clusters
kubectl label no kind-control-plane kubernetes.io/role=master
