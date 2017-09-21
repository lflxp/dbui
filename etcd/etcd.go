package etcd

import (
	//_ "github.com/lflxp/databases/routers"
	//"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"time"
	//"fmt"
	"context"
	"strings"
)

type EtcdUi struct {
	Endpoints 	[]string  //[]string{"localhost:2379"}
	ClientConn 	*clientv3.Client
	Tree 		[]map[string]string
	TopName 	[]string 	//顶级数据库名
}

func (this *EtcdUi) Remove(s []string,de string) []string {
	tmp := []string{}
	for _,key := range s {
		if key != de {
			tmp = append(tmp,key)
		}
	}
	return tmp
}

//获取顶级目录
//TopName目前只支持单数据库 多数据库以后写
func (this *EtcdUi) GetTopic(data []string) []string {
	tmp := data
	for _,key := range data {
		for _,k2 := range data {
			if key != k2 {
				//如果key短值是k2长值得开头(key比k2短，但是key和k2得值是包含关系)
				if strings.HasPrefix(k2,key) && strings.Contains(k2,key) {
					//fmt.Println("##############TOP ",k2,key)
					tmp = this.Remove(tmp,k2)
				}
			}
		}
	}
	this.TopName = []string{"ETCD->"+strings.Join(this.Endpoints,"<-")}
	//将连接作为顶级域名
	this.Tree = append(this.Tree,map[string]string{"name":this.TopName[0],"parentOrg":"null"})
	for _,k := range tmp {
		ttt := map[string]string{}
		ttt["name"] = k
		ttt["value"] = k
		ttt["parentOrg"] = this.TopName[0]
		this.Tree = append(this.Tree,ttt)
	}
	return tmp
}

//判断现有tree集合里面有key没有
func (this *EtcdUi) HasKeyByTree(key string) bool {
	rs := false
	if len(this.Tree) == 0 {
		return false
	}
	for _,k := range this.Tree {
		if value,ok := k["name"]; ok {
			if value == key {
				rs = true
			}
		}
	}
	return rs
}

//根据顶级目录 获取所有子目录
func (this *EtcdUi) GetLastData(key string) {
	last := this.More(key)
	for _,y := range last.Kvs {
		if string(y.Key) != key {
			tmp := map[string]string{}
			//fmt.Println("getLastData",string(y.Key))
			if this.HasChildTree(string(y.Key)) {
				if !this.HasKeyByTree(string(y.Key)) {
					//fmt.Println("has more",string(y.Key),key)
					tmp["name"] = string(y.Key)
					tmp["value"] = string(y.Value)
					tmp["parentOrg"] = key
					this.Tree = append(this.Tree,tmp)
				}
				this.GetLastData(string(y.Key))
			} else {
				if !this.HasKeyByTree(string(y.Key)) {
					//fmt.Println("no more",string(y.Key),string(y.Value),key)
					tmp["name"] = string(y.Key)
					tmp["value"] = string(y.Value)
					tmp["parentOrg"] = key
					this.Tree = append(this.Tree,tmp)
				}
			}
		}
	}
}

func (this *EtcdUi) HasChildTree(key string) bool {
	var status bool
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	resp,err := this.ClientConn.Get(ctx,key,clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	if resp.Count == 1 || resp.Count == 0 {
		status = false
	} else {
		status = true
	}
	return status
}


//more 是底层吗
func (this *EtcdUi) More(data string) *clientv3.GetResponse {
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	resp,err := this.ClientConn.Get(ctx,data,clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	return resp
}

//获取数据
func (this *EtcdUi) FindData(data string) map[string]interface{} {
	defer this.Close()
	TotalRs := map[string]interface{}{}
	result := []map[string]interface{}{}
	this.InitClientConn()
	resp := this.More(data)
	TotalRs["total"] = resp.Count
	for _,key := range resp.Kvs {
		tmp := map[string]interface{}{}
		tmp["id"] = string(key.Key)
		tmp["value"] = string(key.Value)
		tmp["version"] = key.Version
		tmp["lease"] = key.Lease
		tmp["createrevision"] = key.CreateRevision
		tmp["moderevision"] = key.ModRevision
		tmp["memberid"] = resp.Header.MemberId
		tmp["ClusterId"] = resp.Header.ClusterId
		tmp["RaftTerm"] = resp.Header.RaftTerm
		result = append(result,tmp)
	}
	TotalRs["rows"] = result
	return TotalRs
}

//more 是底层吗
func (this *EtcdUi) Count(data string) bool {
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	resp,err := this.ClientConn.Get(ctx,data,clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	return resp.More
}

//get etcd clientV3 conn
func (this *EtcdUi) InitClientConn() {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:this.Endpoints,
		DialTimeout:5*time.Second,
	})
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	this.ClientConn = cli
}

func (this *EtcdUi) Close() {
	this.ClientConn.Close()
}

//endpoint []string{"localhost:2379"}
func (this *EtcdUi) GetAllDatas() []string {
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	resp,err := this.ClientConn.Get(ctx,"",clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	//fmt.Println(resp.More)
	rs := []string{}
	for _,ev := range resp.Kvs {
		//fmt.Println(string(ev.Key))
		//fmt.Println(fmt.Printf("%d %s %s\n",ev.Lease,ev.Key,ev.Value))
		//fmt.Println(ev.String())
		//rs[string(ev.Key)] = string(ev.Value)
		rs = append(rs,string(ev.Key))
	}
	return rs
}

//快速获取tree所有信息
func (this *EtcdUi) GetAllTreeRelate() {
	this.InitClientConn()
	for _,key := range this.GetTopic(this.GetAllDatas()) {
		this.GetLastData(key)
	}
}

//递归函数，获得树形结构的关系信息 tree-view
func (this *EtcdUi) GetTreeRelate(top []string,all []map[string]string) string {
	result := []string{}
	for _,y := range top {
		result = append(result,"{text:'"+strings.Split(y,"/")[len(strings.Split(y,"/"))-1]+"'")
		if this.HasChild(y,all) {
			result = append(result,"selectable:true,multiSelect:false,state:{expanded:false,disabled:false},href:'#',ids:'"+y+"','nodes':["+this.GetTreeRelate(this.ForeignKeys(y,all),all)+"]}")
		} else {
			result = append(result,"icon:'glyphicon glyphicon-user',selectable:true,href:'#',ids:'"+y+"'}")
		}
	}
	return strings.Join(result,",")
}

//判断是否还有子机构
func (this *EtcdUi) HasChild(id string,data []map[string]string) bool {
	ok := false
	for _,y := range data {
		if y["parentOrg"] == id {
			ok = true
		}
	}
	return ok
}

//获取所有上级机构为key的子机构（第二层，用于下面的递归）
func (this *EtcdUi) ForeignKeys(key string,data []map[string]string) []string {
	res := []string{}
	for _,y := range data {
		if y["parentOrg"] == key {
			res = append(res,y["name"])
		}
	}
	return res
}


//根据顶级机构和所有数据进行递归 得到树形结构的json字符串
//获取所有tree table最终数据
func (this *EtcdUi) GetTreeByString() string {
	defer this.Close()
	this.GetAllTreeRelate()
	//return "["+this.GetTreeRelate(this.GetTopic(this.GetAllDatas()),this.Tree)+"]"
	return "["+this.GetTreeRelate(this.TopName,this.Tree)+"]"
}