package main

import (
	"fmt"
	"gather/toolkitcl/demon/casbin/adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"log"
	"os"
)

func main()  {
	currentPath, _ := os.Getwd()
	fmt.Println("getTmpDir(当前系统临时目录) = ", currentPath)
	// CodeModelAndAdapter("./model/keymatch2_policy.csv")
	// FileModelAndAdapter("./model/keymatch2_model.conf", "./model/keymatch2_policy.csv")
	// StringModelAndAdapter("./model/keymatch2_policy.csv")
	StringModelAndStringAdapter()
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

func CodeModelAndAdapter(policyPath string)  {
	m := model.NewModel()
	m.AddDef("r","r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")

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

func StringModelAndAdapter(policyPath string)  {
	modelDefinition :=
		`
		[request_definition]
        r = sub, obj, act

        [policy_definition]
        p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
		`
	m, _ := model.NewModelFromString(modelDefinition)

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

func StringModelAndStringAdapter()  {
	modelDefinition :=
		`
		[request_definition]
        r = sub, obj, act

        [policy_definition]
        p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
		`
	policyDefinition :=
		`
		p, alice, /alice_data/*, GET
		p, alice, /alice_data/resource1, GET

		p, bob, /alice_data/resource2, GET
		p, bob, /bob_data/*, POST

		p, cathy, /cathy_data, (GET)|(POST)
        `
	m, _ := model.NewModelFromString(modelDefinition)

	a := adapter.NewAdapter(policyDefinition)

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





