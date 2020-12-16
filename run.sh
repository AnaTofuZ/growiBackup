DAY=`date "+%Y-%m-%d"`
REVJSON=revision_back_${DAY}.json
cd $HOME/backup
docker exec growi_mongo_1 mongoexport -d growi -c revisions --pretty --jsonArray --out ${REVJSON}
docker cp growi_mongo_1:${REVJSON} .
docker exec growi_mongo_1 rm ${REVJSON}
./growibackup ${REVJSON} Growi
ls *.json | grep -v ${REVJSON} | xargs rm

cd Growi
hg add .
hg commit -u 'autobackup' -m "backup ${DAY}"
hg push
