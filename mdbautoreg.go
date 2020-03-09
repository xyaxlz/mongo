package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mgo.v2"
	"net"
	"strings"
	//	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os/exec"
	//	"gopkg.in/mgo.v2/bson"
)

type Mongo_Result_Repl_Status struct {
	Set     string                 `bson:"set"`
	Members []Mongo_Replica_Member `bson:"members"`
}

type Mongo_Replica_Member struct {
	Id        int    `bson:"_id"`
	Name      string `bson:"name"`
	Health    int    `bson:"health"`
	State     bool   `bson:"state"`
	StateStr  string `bson:"stateStr"`
	SyncingTo string `bson:"syncingTo"`
	Self      bool   `bson:"self"`
}

type MdbReg struct {
	Ip    string
	Port  int
	MIp   string
	MPort int
	Flag  int
}

type Mongo_List_DB struct {
	Databases []DB_Member `bson:"databases"`
	TotalSize int         `bson:"totalSize"`
}

type DB_Member struct {
	Name string `bson:"name"`
}

var (
	regUserName = "dbreg"
	regPassword = "xxxx"
	regIp       = "1xxxxx"
	regDbName   = "db_asset"
	regPort     = "3388"
)
var DB *sql.DB

func checkErr(err error) {

	if err != nil {
		fmt.Println(err)
		panic(err)

	}

}

func execShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	checkErr(err)
	return out.String(), err

}

func insertDB(hostIp string, hostPort string, mIp string, mPort string, flag string, dbsName string) {
	//生成数据库连接地址
	regPath := strings.Join([]string{regUserName, ":", regPassword, "@tcp(", regIp, ":", regPort, ")/", regDbName, "?charset=utf8"}, "")
	//连接数据库
	regDB, _ := sql.Open("mysql", regPath)
	//判断数据库是否连接成功，如果失败打印日志到中断，并退出
	if err := regDB.Ping(); err != nil {
		fmt.Printf("open %s %s database fail\n", regIp, regPort)
		return
	}

	//prepare insert 语句
	stmt, err := regDB.Prepare(`replace into mdb_asset (ip, port, mip, mport, flag, dbs) values (?, ?, ?, ?, ?, ?)`)
	checkErr(err)

	_, err = stmt.Exec(hostIp, hostPort, mIp, mPort, flag, dbsName)
	checkErr(err)
	fmt.Printf("replace into mdb_asset (ip, port, mip, mport, flag, dbs) values (%s, %s, %s, %s, %s, %s) \n", hostIp, hostPort, mIp, mPort, flag, dbsName)
}

func getMongos(port string) (dbsRegArr [][]string) {
	var dbsArr []string
	var dbsStr string
	//	var returnNil [][]string
	var primaryName string

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err)  // 这里的err其实就是panic传入的内容
			fmt.Println(port) // 这里的err其实就是panic传入的内容
			//return returnNil

		}
	}()

	fmt.Println(port)
	session, err := mgo.Dial("127.0.0.1:" + port)
	session.SetMode(mgo.Monotonic, true)

	err = session.DB("admin").Login("root", "xxx")
	if err != nil {
		//fmt.Println(err)
		return dbsRegArr
	}

	result1 := Mongo_List_DB{}
	err = session.DB("admin").Run("listDatabases", &result1)
	for _, v := range result1.Databases {
		if v.Name != "admin" && v.Name != "local" {
			dbsArr = append(dbsArr, v.Name)
		}
	}

	//fmt.Println(dbsArr)
	dbsStr = strings.Join(dbsArr, " ")

	result := Mongo_Result_Repl_Status{}
	err = session.DB("admin").Run("replSetGetStatus", &result)
	//fmt.Println(result)
	if err != nil {
		fmt.Println("------连接数据库失败------------")
		panic(err)
	}

	//fmt.Println(result.Members)
	for _, v := range result.Members {
		if v.StateStr == "PRIMARY" {
			primaryName = v.Name
		}
	}
	//fmt.Println(primaryName)
	for _, v := range result.Members {
		if v.StateStr != "PRIMARY" {
			ns1, _ := net.LookupHost(strings.Split(v.Name, ":")[0])
			ns2, _ := net.LookupHost(strings.Split(primaryName, ":")[0])
			//dbsRegArr = append(dbsRegArr, []string{strings.Split(v.Name, ":")[0], strings.Split(v.Name, ":")[1], strings.Split(primaryName, ":")[0], strings.Split(primaryName, ":")[1], "1", dbsStr})
			dbsRegArr = append(dbsRegArr, []string{ns1[0], strings.Split(v.Name, ":")[1], ns2[0], strings.Split(primaryName, ":")[1], "1", dbsStr})
		} else {

			ns2, _ := net.LookupHost(strings.Split(primaryName, ":")[0])
			//dbsRegArr = append(dbsRegArr, []string{strings.Split(v.Name, ":")[0], strings.Split(v.Name, ":")[1], "", "", "0", dbsStr})
			dbsRegArr = append(dbsRegArr, []string{ns2[0], strings.Split(v.Name, ":")[1], "", "", "0", dbsStr})
		}
	}
	return dbsRegArr

}

func main() {
	mongoPorts, err := execShell("netstat  -lntp |grep -w mongod|awk '{print $4}'|awk -F ':' '{print $NF}' |sort|uniq")
	checkErr(err)
	for _, port := range strings.Fields(mongoPorts) {
		//fmt.Println(port)
		for _, regDB := range getMongos(port) {
			//fmt.Println(regDB)
			insertDB(regDB[0], regDB[1], regDB[2], regDB[3], regDB[4], regDB[5])
		}
	}

}
