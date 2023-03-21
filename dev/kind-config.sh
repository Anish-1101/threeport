#!/usr/bin/env bash

CODE_DIR=$(pwd)

cat << EOF > dev/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: threeport-dev
kubeadmConfigPatches:
- |-
  kind: ClusterConfiguration
  # configure controller-manager bind address
  controllerManager:
    extraArgs:
      bind-address: 0.0.0.0
  # configure etcd metrics listen address
  etcd:
    local:
      extraArgs:
        listen-metrics-urls: http://0.0.0.0:2381
  # configure scheduler bind address
  scheduler:
    extraArgs:
      bind-address: 0.0.0.0
- |-
  kind: KubeProxyConfiguration
  # configure proxy metrics bind address
  metricsBindAddress: 0.0.0.0
nodes:
- role: control-plane
  kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
  extraPortMappings:
    - containerPort: 80
      hostPort: 80
      protocol: TCP
    - containerPort: 443
      hostPort: 443
      protocol: TCP
- role: worker
  extraMounts:
    - hostPath: $CODE_DIR
      containerPath: /threeport-rest-api
- role: worker
  extraMounts:
    - hostPath: $CODE_DIR
      containerPath: /threeport-rest-api
- role: worker
  extraMounts:
    - hostPath: $CODE_DIR
      containerPath: /threeport-rest-api
EOF

#cat << EOF > dev/kind-config.yaml
#kind: Cluster
#apiVersion: kind.x-k8s.io/v1alpha4
#name: threeport-dev
#nodes:
#- role: control-plane
#- role: worker
#  extraMounts:
#    - hostPath: $CODE_DIR
#      containerPath: /threeport-rest-api
#- role: worker
#  extraMounts:
#    - hostPath: $CODE_DIR
#      containerPath: /threeport-rest-api
#- role: worker
#  extraMounts:
#    - hostPath: $CODE_DIR
#      containerPath: /threeport-rest-api
#EOF

echo "KinD config file generated"

exit 0
