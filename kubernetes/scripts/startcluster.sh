#!/bin/zsh

NODE_NAME_PREFIX=ad
NODE_CPUS=2
NODE_MEM=4g
NODE_DISK=15g

NODE_COUNT=3

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
  multipass Provision --name "$NODE_NAME" --cpus $NODE_CPUS --mem $NODE_MEM --disk $NODE_DISK

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
  multipass exec "$NODE_NAME" -- bash -c "
  sudo yq -i eval '.clusters[].cluster.server |= sub(\\\"127.0.0.1\\\", env(NODE_IP))
   | .contexts[].name = \\\"$NODE_NAME\\\" | .current-context = \\\"$NODE_NAME\\\" |
   .clusters[].name=\\\"$NODE_NAME\\\" | .contexts[].context.cluster=\\\"$NODE_NAME\\\"
   | .users[].name=\\\"$NODE_NAME\\\" | .contexts[].context.user=\\\"$NODE_NAME\\\"' \
   $KUBECONFIG_NAME"

  multipass exec "$NODE_NAME" -- bash -c "sudo chmod 777 $KUBECONFIG_NAME"

  if [ "$node_id" -eq 1 ]
  then
    multipass transfer "$NODE_NAME":/home/ubuntu/"$KUBECONFIG_NAME" "$HOME"/.kube/"$KUBECONFIG_NAME"
  else
    echo "Transferring $KUBECONFIG_NAME to $NODE_1_NAME"
    multipass transfer "$NODE_NAME":/home/ubuntu/"$KUBECONFIG_NAME" "$HOME"/.kube/"$KUBECONFIG_NAME"
    multipass transfer "$HOME"/.kube/"$KUBECONFIG_NAME" $NODE_1_NAME:/home/ubuntu/"$KUBECONFIG_NAME"
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

# Installing krew on local machine
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
echo export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH" >> ~/.zshrc
echo $PATH

# Installing konfig for kubectl config management
kubectl krew install konfig


for kf in $NODE_RANGE
do
  node_name="$NODE_NAME_PREFIX$kf"
  kf_name=$HOME/.kube/kubeconfig.$node_name
  echo "Merging konfig $kf_name"
  kubectl config delete-cluster "$node_name"
  kubectl config delete-context "$node_name"
  kubectl config delete-user "$node_name"
  kubectl konfig import -s "$kf_name"
done

username=$(cat ~/.credentials/arya.pat | grep "username" | awk '{print $2}')
password=$(cat ~/.credentials/arya.pat | grep "password" | awk '{print $2}')

# Deploying our helm chart
for n in $NODE_RANGE
do

  node_name="$NODE_NAME_PREFIX$n"
  kubectl config use-context $node_name
  # Adding an access token secret

  kubectl delete secret regcred
  kubectl create secret generic regcred \
  --from-file=.dockerconfigjson=$HOME/.docker/config.json \
  --type=kubernetes.io/dockerconfigjson

  if [ $n -eq 1 ]
  then
    kubectl label nodes $node_name aryaRole="orchestrator"
  fi
  join_addrs=""

  node_ip=$(multipass info $node_name | grep IPv4 | awk '{print $2}')
  for nx in $NODE_RANGE
  do
    if [ $nx -ne $n ]
    then
      nx_name="$NODE_NAME_PREFIX$nx"
      join_addr=$(multipass info $nx_name | grep IPv4 | awk '{print $2}')
      join_addrs+="$join_addr\,"
    fi
  done
  join_addrs=${join_addrs%?}
  echo $join_addrs
  echo $node_ip
  helm uninstall aryacore
  helm install --set cockroachdb.clusterInitHost=ad1,cockroachdb.nodeIP=$node_ip,cockroachdb.join=$join_addrs\
  aryacore ./aryacore
done