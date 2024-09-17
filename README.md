scripts
migrate up cmd : 'go run main.go migrate-up'
migrate down cmd : 'go run main.go migrate-down{version}' , 'docker exec -it {container_id} go run main.go migrate-down{version}'


env example
SERVER_PORT=8080
DSN=root:root@tcp(localhost:3306)/kab-phone?charset=utf8mb4&parseTime=True
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
JWT_SECRET=kab-phone
ADMIN_PASSWORD=4dm1n
MAX_IMAGE_SIZE=5