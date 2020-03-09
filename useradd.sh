while read line 
do
	echo $line
	/data/server/mongodb-linux-x86_64-rhel62-3.2.13/bin/mongo ${line}/admin  -uroot -p'xxxx' < useradd.js
	#/data/server/mongodb-linux-x86_64-rhel62-3.2.13/bin/mongo ${line}/admin  -uroot -p'ss' --eval "printjson(db.getCollectionNames())"
#done < master
done < master2
