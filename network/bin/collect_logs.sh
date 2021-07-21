
for container_name in $(docker ps --format "{{.Names}}");
    do 

        sudo cp $(docker inspect --format='{{.LogPath}}' $container_name) ./logs/"$container_name.log"
        #echo $(docker logs $(docker ps -aqf "name=$container_name")) > ./logs/"$container_name.log";
    done