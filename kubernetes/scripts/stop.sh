NODE_COUNT=4
NODE_RANGE=$(seq 1 $NODE_COUNT)

for id in $NODE_RANGE
do
  name="ad$id"
  multipass stop $name
  multipass delete $name
done

multipass purge