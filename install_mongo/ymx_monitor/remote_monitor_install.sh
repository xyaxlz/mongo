for i in `seq 5`
do 
	echo $i
	ssh mdb${i}v.infra.bjza.pdtv.it "mkdir -p /data/install/monitor"
	scp monitor_install.sh  mdb${i}v.infra.bjza.pdtv.it:/data/install/monitor
	scp mongodb.conf  mdb${i}v.infra.bjza.pdtv.it:/data/install/monitor
	ssh  mdb${i}v.infra.bjza.pdtv.it "cd /data/install/monitor;sh monitor_install.sh"
	#scp mongodb.conf  mdb${i}v.infra.bjza.pdtv.it:/etc/zabbix/zabbix_agentd.d/
	#a=`ssh -n  mdb${i}v.infra.bjza.pdtv.it 'awk  "/port:/{print $1}" /etc/mongod.conf'|awk '{print $2}'`
	#a=`ssh -n  mdb${i}v.infra.bjac.pdtv.it 'awk  "/port:/{print $1}" /etc/mongod.conf'`
	#echo $a
	#ssh mdb${i}v.infra.bjza.pdtv.it "sed -i 's/7090/$a/g' /etc/zabbix/zabbix_agentd.d/mongodb.conf   "
	#ssh mdb${i}v.infra.bjza.pdtv.it "/etc/init.d/zabbix-agent restart"
done
