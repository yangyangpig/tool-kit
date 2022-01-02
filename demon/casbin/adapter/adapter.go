package adapter

import (
	"bytes"
	"errors"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/util"
	"strings"
)

type Adapter struct {
	Policy string
}

func NewAdapter(policy string) *Adapter {
	return &Adapter{
		Policy: policy,
	}
}

func (a *Adapter) LoadPolicy(model model.Model) error {
	if a.Policy == "" {
		return errors.New("invalid policy string, policy cannot be empty")
	}
	splitStr := strings.TrimSpace(a.Policy)
	strs := strings.Split(splitStr, "\n")
	for _, str := range strs {

		if str == "" {
			continue
		}
		persist.LoadPolicyLine(str, model)
	}
	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	var tmp bytes.Buffer
	// 获取policy的内容,并拼接成policy的model格式
	for ptype, val := range model["p"] {
		for _, rule := range val.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}

	// 获取rbac的g类型内容，并且拼接成policy的model格式
	for gtype, val := range model["g"] {
		for _, rule := range val.Policy {
			tmp.WriteString(gtype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}
	a.Policy = strings.TrimRight(tmp.String(), "\n")
	return nil
}

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	a.Policy = ""
	return nil
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
