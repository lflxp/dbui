package etcd

import (
	"testing"
	"strings"
	"fmt"
)

func TestEtcdUi_GetTopic(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	st.InitClientConn()
	top := st.GetTopic(st.GetAllDatas())
	t.Log(strings.Join(top,":"))
}

func TestEtcdUi_GetLastData(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	st.InitClientConn()
	for _,key := range st.GetTopic(st.GetAllDatas()) {
		st.GetLastData(key)
	}
	fmt.Println("############")
	for _,k := range st.Tree {
		fmt.Println(fmt.Sprintf("%s %s %s",k["name"],k["value"],k["parentOrg"]))
	}
}

func TestEtcdUi_More(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	st.InitClientConn()
	rs := st.More("foo1123")
	fmt.Println(fmt.Sprintf("%d",rs.Count))
	for _,k := range rs.Kvs {
		fmt.Println(string(k.Key),string(k.Value))
	}
}

func TestEtcdUi_HasChildTree(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	st.InitClientConn()
	fmt.Println(st.HasChildTree("foo11233"))
}

func TestEtcdUi_GetAllTreeRelate(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	st.GetAllTreeRelate()
	for _,k := range st.Tree {
		fmt.Println(fmt.Sprintf("%s %s %s",k["name"],k["value"],k["parentOrg"]))
	}
}

func TestEtcdUi_GetTreeByString(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	rs := st.GetTreeByString()
	t.Log(rs)
}

func TestEtcdUi_Add(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	err := st.Add("/lxp","good")
	err = st.Add("/lxp/good","good")
	err = st.Add("/lxp/good1","good")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("ok")
}

func TestEtcdUi_Delete(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	err := st.Delete("/lxp/good")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("ok")
}

func TestEtcdUi_DeleteAll(t *testing.T) {
	st := EtcdUi{Endpoints:[]string{"localhost:2379"}}
	err := st.DeleteAll("/lxp/good")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("ok")
}