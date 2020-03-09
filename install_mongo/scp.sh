for i in {4,5,6,7,8,9}
do 
	echo $i
scp 030.host.default  mdb${i}v.infra.bjac.pdtv.it:/home/server_config/iptables_rules/
ssh mdb${i}v.infra.bjac.pdtv.it "/etc/init.d/iptables restart"
ssh mdb${i}v.infra.bjac.pdtv.it "/etc/init.d/zabbix-agent restart"
done
