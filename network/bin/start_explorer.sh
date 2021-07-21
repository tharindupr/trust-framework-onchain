echo '================ Stoping previous Explorere Instaces ================'
docker-compose -f ./../blockchain-explorer/docker-compose.yaml down -v


echo    '================ Starting the Docker Instances ================'
docker-compose -f ./../blockchain-explorer/docker-compose.yaml up -d
