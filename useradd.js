use admin
db.createUser(
	  {
		      user: "monitor_dba",
		      pwd: "xxxxx",
		      roles: [ { role: "root", db: "admin" } ]
		    })

