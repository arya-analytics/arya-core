#!/bin/zsh

NODE_NAME_PREFIX=ad
NODE_CPUS=2
NODE_MEM=4g
NODE_DISK=15g

NODE_COUNT=4

NODE_1_NAME=$NODE_NAME_PREFIX
NODE_1_NAME+=1
NODE_RANGE=$(seq 1 $NODE_COUNT)

echo "Initializing Arya Development Cluster with $NODE_COUNT nodes"

for node_id in $NODE_RANGE
do
  CIDR_OFFSET=$((node_id*2))
  POD_CIDR_BASE=$((43+CIDR_OFFSET))
  SERVICE_CIDR_BASE=$((POD_CIDR_BASE+1))
  POD_CIDR=10.
  POD_CIDR+=$POD_CIDR_BASE
  POD_CIDR+=.0.0/16
  SERVICE_CIDR=10.
  SERVICE_CIDR+=$SERVICE_CIDR_BASE
  SERVICE_CIDR+=.0.0/16

  NODE_NAME=$NODE_NAME_PREFIX
  NODE_NAME+=$node_id

  # Starting multipass VM
  multipass launch --name "$NODE_NAME" --cpus $NODE_CPUS --mem $NODE_MEM --disk $NODE_DISK

  # shellcheck disable=SC2086
  NODE_IP=$(multipass info $NODE_NAME | grep IPv4 | awk '{print $2}')

  echo "Node started at $NODE_IP"
  echo "Starting k3s cluster with pod cidr $POD_CIDR and service cidr $SERVICE_CIDR"

  # Installing k3s
  multipass exec "$NODE_NAME" -- bash -c "curl -sfL https://get.k3s.io |
  INSTALL_K3S_EXEC=\"--cluster-cidr $POD_CIDR  --service-cidr $SERVICE_CIDR \
  --write-kubeconfig-mode 777 \" sh -s - "

  # Labeling a node a worker
  multipass exec "$NODE_NAME" -- bash -c "kubectl label node $NODE_NAME node-role.kubernetes.io/worker=worker"

  # Installing yq for yaml editing
  # shellcheck disable=SC2086
  multipass exec $NODE_NAME -- bash -c "sudo snap install yq"

  # Setting up kubeconfig
  echo "Setting up kubeconfig"
  KUBECONFIG_NAME=kubeconfig.$NODE_NAME
  multipass exec "$NODE_NAME" -- bash -c "sudo cp /etc/rancher/k3s/k3s.yaml $KUBECONFIG_NAME"
  multipass exec "$NODE_NAME" -- bash -c "sudo bash -c \"echo NODE_IP=$NODE_IP >> \
  /etc/environment\""
  multipass exec "$NODE_NAME" -- bash -c "echo $NODE_IP"
  # shellcheck disable=SC2086
  multipass exec $NODE_NAME -- bash -c "
  sudo yq -i eval '.clusters[].cluster.server |= sub(\\\"127.0.0.1\\\", env(NODE_IP))
   | .contexts[].name = \\\"$NODE_NAME\\\" | .current-context = \\\"$NODE_NAME\\\" |
   .clusters[].name=\\\"$NODE_NAME\\\" | .contexts[].context.cluster=\\\"$NODE_NAME\\\"
   | .users[].name=\\\"$NODE_NAME\\\" | .contexts[].context.user=\\\"$NODE_NAME\\\"' \
   $KUBECONFIG_NAME"

  multipass exec "$NODE_NAME" -- bash -c "sudo chmod 777 $KUBECONFIG_NAME"

  if [ "$node_id" -eq 1 ]
  then
    echo "Installing Submariner CLI on node $NODE_NAME"
    multipass exec "$NODE_NAME" -- bash -c "
      curl -Ls https://get.submariner.io | bash
      export PATH=\$PATH:~/.local/bin
      echo export PATH=\$PATH:~/.local/bin >> ~/.profile
      subctl deploy-broker --kubeconfig $KUBECONFIG_NAME"

  else
    echo "Transferring $KUBECONFIG_NAME to $NODE_1_NAME"
    multipass transfer "$NODE_NAME":/home/ubuntu/"$KUBECONFIG_NAME" ~/"$KUBECONFIG_NAME"
    multipass transfer ~/"$KUBECONFIG_NAME" $NODE_1_NAME:/home/ubuntu/"$KUBECONFIG_NAME"
    rm ~/kubeconfig."$NODE_NAME"
  fi

  echo "All done setting up cluster $NODE_NAME"
done

KUBECONFIGS=""
for kf in $NODE_RANGE
do
  KUBECONFIGS+="kubeconfig.$NODE_NAME_PREFIX$kf"
  if [ "$kf" -ne $NODE_COUNT ]
  then
    KUBECONFIGS+=":"
  fi
done

echo "$KUBECONFIGS"
multipass exec $NODE_1_NAME -- bash -c "
  sudo bash -c \"echo KUBECONFIG=$KUBECONFIGS >> /etc/environment\"
"

for node_idx in $NODE_RANGE
do
  node_name="$NODE_NAME_PREFIX$node_idx"
  kubeconfig_name="\$HOME/kubeconfig.$node_name"
  echo "Joining cluster $node_name to the submariner broker on $NODE_1_NAME"
  multipass exec $NODE_1_NAME -- bash -c "
  export PATH=\$PATH:~/.local/bin
  echo export PATH=\$PATH:~/.local/bin >> ~/.profile
  subctl join --kubeconfig $kubeconfig_name broker-info.subm --clusterid \
  $node_name --natt=false"
done
