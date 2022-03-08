go run main.go demo config --node-addr=localhost:26260 --node=1 && go run main.go demo config --node-addr=localhost:26260 --node=2 &&go run main.go demo config --node-addr=localhost:26260 --node=3

go run main.go demo write --node-addr=localhost:26261 --node=1

go run main.go demo write --node-addr=localhost:26259 --node=2

go run main.go demo write --node-addr=localhost:26260 --node=3

docker container rm minio1 minio2 minio3 roach1 roach2 roach3 && docker volume rm docker_minio1-data docker_minio2-data docker_minio3-data docker_roach1-data docker_roach2-data docker_roach3-data