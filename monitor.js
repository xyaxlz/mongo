use admin
db.dropUser('monitor')
db.createUser({
	user: "monitor",
	pwd: "xxx",
	roles:[
		            {"db": "admin", "role": "clusterMonitor"},
		            {"db": "admin", "role": "read"}
		        ]
})
