#!/bin/bash

backup_path="/home/docker/mysql/backup"
docker="/usr/local/bin/docker-compose"
# 目录是否存在
if [ ! -d "$backup_path" ]; then
   mkdir -p $backup_path
fi
# 指定目录
cd /home/docker
# 导出数据库
docker exec mysql-master  bash -c "mysqldump -h localhost -uroot -p\$MYSQL_ROOT_PASSWORD car > /var/backup/$(date +%Y%m%d).sql"
# 删除超过10天的数据
rm -f $backup_path/$(date -d -10day +%Y%m%d).sql
