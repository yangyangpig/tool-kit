package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"log"
	"os"
)

func main()  {
	currentPath, _ := os.Getwd()
	fmt.Println("getTmpDir(当前系统临时目录) = ", currentPath)
	CodeModelANdRedisAdapter("./model/keymatch2_policy.csv")
	FileModelAndAdapter("./model/keymatch2_model.conf", "./model/keymatch2_policy.csv")
}

func FileModelAndAdapter(modelPath, policyPath string)  {
	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatalf("new enforcer happen error %+v", err)
		return
	}

	authResp, err := e.BatchEnforce([][]interface{}{{"alice","/alice_data/v2","GET"}, {"alice","/alice_data/resource1","GET"}, {"cathy","/cathy_data","(GET)|(POST)"}})
	if err != nil {
		log.Fatalf("enforce happen error %+v", err)
		return
	}

	fmt.Printf("authResp : %v", authResp)
}

func CodeModelANdRedisAdapter(policyPath string)  {
	m := model.NewModel()
	m.AddDef("r","r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.obj == p.obj && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)")

	a := fileadapter.NewAdapter(policyPath)

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("new enforcer happen error: %+v", err)
		return
	}
	authResp, err := e.BatchEnforce([][]interface{}{{"alice","/alice_data/v2","GET"}, {"alice","/alice_data/resource1","GET"}, {"cathy","/cathy_data","(GET)|(POST)"}})
	if err != nil {
		log.Fatalf("enforce happen error %+v", err)
		return
	}

	fmt.Printf("authResp : %v", authResp)
}

func redisAdapter()  {
	
}





