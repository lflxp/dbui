package etcd

import (
	"fmt"
	//_ "github.com/lflxp/databases/routers"
	//"github.com/astaxie/beego"
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

type EtcdUi struct {
	Endpoints    []string //[]string{"localhost:2379"}
	ClientConn   *clientv3.Client
	ClientConnV2 client.Client
	Tree         []map[string]string
	TopName      []string //顶级数据库名
	Version      string   //2 or 3
}

func (this *EtcdUi) Remove(s []string, de string) []string {
	tmp := []string{}
	for _, key := range s {
		if key != de {
			tmp = append(tmp, key)
		}
	}
	return tmp
}

//获取顶级目录
//TopName目前只支持单数据库 多数据库以后写
func (this *EtcdUi) GetTopic(data []string) []string {
	tmp := data
	for _, key := range data {
		for _, k2 := range data {
			if key != k2 {
				//如果key短值是k2长值得开头(key比k2短，但是key和k2得值是包含关系)
				if strings.HasPrefix(k2, key) && strings.Contains(k2, key) {
					//fmt.Println("##############TOP ",k2,key)
					tmp = this.Remove(tmp, k2)
				}
			}
		}
	}
	this.TopName = []string{"ETCD->" + strings.Join(this.Endpoints, "<-")}
	//将连接作为顶级域名
	this.Tree = append(this.Tree, map[string]string{"name": this.TopName[0], "parentOrg": "null"})
	for _, k := range tmp {
		ttt := map[string]string{}
		ttt["name"] = k
		ttt["value"] = k
		ttt["parentOrg"] = this.TopName[0]
		this.Tree = append(this.Tree, ttt)
	}
	return tmp
}

//判断现有tree集合里面有key没有
func (this *EtcdUi) HasKeyByTree(key string) bool {
	rs := false
	if len(this.Tree) == 0 {
		return false
	}
	for _, k := range this.Tree {
		if value, ok := k["name"]; ok {
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
	for _, y := range last.Kvs {
		if string(y.Key) != key {
			tmp := map[string]string{}
			//fmt.Println("getLastData",string(y.Key))
			if this.HasChildTree(string(y.Key)) {
				if !this.HasKeyByTree(string(y.Key)) {
					//fmt.Println("has more",string(y.Key),key)
					tmp["name"] = string(y.Key)
					tmp["value"] = string(y.Value) 
					tmp["ttl"] = fmt.Sprintf("%d", y.Lease)
					tmp["version"] = fmt.Sprintf("%d", y.Version)
					tmp["parentOrg"] = key
					this.Tree = append(this.Tree, tmp)
				}
				this.GetLastData(string(y.Key))
			} else {
				if !this.HasKeyByTree(string(y.Key)) {
					//fmt.Println("no more",string(y.Key),string(y.Value),key)
					st1 := string(y.Value)
					st2 := strings.Split(st1, "::")
					tmp["name"] = string(y.Key)
					tmp["value"] = st1
					if len(st2) == 7 {
						tmp["did"] = st2[0]
						tmp["dimage"] = st2[1]
						tmp["dcommand"] = st2[2]
						tmp["dstate"] = st2[3]
						tmp["dStatus"] = st2[4]
						tmp["dports"] = st2[5]
						tmp["dname"] = st2[6]
					}
					tmp["ttl"] = fmt.Sprintf("%d", y.Lease)
					tmp["version"] = fmt.Sprintf("%d", y.Version)
					tmp["parentOrg"] = key
					this.Tree = append(this.Tree, tmp)
				}
			}
		}
	}
}

func (this *EtcdUi) HasChildTree(key string) bool {
	var status bool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := this.ClientConn.Get(ctx, key, clientv3.WithPrefix())
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := this.ClientConn.Get(ctx, data, clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	return resp
}

//more 是底层吗
func (this *EtcdUi) Get(data string) *clientv3.GetResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := this.ClientConn.Get(ctx, data)
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
	for _, key := range resp.Kvs {
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
		tmp["op"] = fmt.Sprintf("<a href=\"#passwd\" data-toggle=\"modal\" id=\"install\" class=\"btn btn-success btn-sm\"><i class=\"glyphicon glyphicon-wrench\"></i>修改</a><button onclick=\"Delete('%s')\" class=\"btn btn-danger btn-sm\"><i class=\"glyphicon glyphicon-remove\"></i> 删除 </button>", string(key.Key))
		result = append(result, tmp)
	}
	TotalRs["rows"] = result
	return TotalRs
}

//more 是底层吗
func (this *EtcdUi) Count(data string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := this.ClientConn.Get(ctx, data, clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	return resp.More
}

//get etcd clientV3 conn
func (this *EtcdUi) InitClientConn() {
	if this.Version != "2" {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   this.Endpoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			//fmt.Println(err.Error())
			panic(err)
		}
		this.ClientConn = cli
	} else {
		cfg := client.Config{
			Endpoints:               this.Endpoints,
			Transport:               client.DefaultTransport,
			HeaderTimeoutPerRequest: 5 * time.Second,
		}
		c, err := client.New(cfg)
		if err != nil {
			panic(err)
		}
		this.ClientConnV2 = c
		kapi := client.NewKeysAPI(c)
		resp, err := kapi.Get(context.Background(), "/asd", &client.GetOptions{Recursive: true})
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Node.Key, resp.Node.Value)
	}
}

func (this *EtcdUi) Close() {
	this.ClientConn.Close()
}

//endpoint []string{"localhost:2379"}
func (this *EtcdUi) GetAllDatas() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := this.ClientConn.Get(ctx, "", clientv3.WithPrefix())
	cancel()
	if err != nil {
		//fmt.Println(err.Error())
		panic(err)
	}
	//fmt.Println(resp.More)
	rs := []string{}
	for _, ev := range resp.Kvs {
		//fmt.Println(string(ev.Key))
		//fmt.Println(fmt.Printf("%d %s %s\n",ev.Lease,ev.Key,ev.Value))
		//fmt.Println(ev.String())
		//rs[string(ev.Key)] = string(ev.Value)
		rs = append(rs, string(ev.Key))
	}
	return rs
}

//快速获取tree所有信息
func (this *EtcdUi) GetAllTreeRelate() {
	this.InitClientConn()
	for _, key := range this.GetTopic(this.GetAllDatas()) {
		this.GetLastData(key)
	}
}

//递归函数，获得树形结构的关系信息 tree-view
func (this *EtcdUi) GetTreeRelate(top []string, all []map[string]string) string {
	result := []string{}
	for _, y := range top {
		result = append(result, "{text:'"+strings.Split(y, "/")[len(strings.Split(y, "/"))-1]+"'")
		if this.HasChild(y, all) {
			result = append(result, "selectable:true,multiSelect:false,state:{expanded:false,disabled:false},href:'#',ids:'"+y+"','nodes':["+this.GetTreeRelate(this.ForeignKeys(y, all), all)+"]}")
		} else {
			result = append(result, "icon:'glyphicon glyphicon-list-alt',selectable:true,href:'#',ids:'"+y+"'}")
		}
	}
	return strings.Join(result, ",")
}

//判断是否还有子机构
func (this *EtcdUi) HasChild(id string, data []map[string]string) bool {
	ok := false
	for _, y := range data {
		if y["parentOrg"] == id {
			ok = true
		}
	}
	return ok
}

//获取所有上级机构为key的子机构（第二层，用于下面的递归）
func (this *EtcdUi) ForeignKeys(key string, data []map[string]string) []string {
	res := []string{}
	for _, y := range data {
		if y["parentOrg"] == key {
			// fmt.Println("HasChild",y["parentOrg"],key,y["name"])
			res = append(res, y["name"])
		}
	}
	return res
}

//根据顶级机构和所有数据进行递归 得到树形结构的json字符串
//获取所有tree table最终数据
func (this *EtcdUi) GetTreeByString() string {
	if this.ScannerPort(this.Endpoints[0]) {
		defer this.Close()
		this.GetAllTreeRelate()
		//return "["+this.GetTreeRelate(this.GetTopic(this.GetAllDatas()),this.Tree)+"]"
		return "[" + this.GetTreeRelate(this.TopName, this.Tree) + "]"
	}
	return fmt.Sprintf("%s Unreachable", this.Endpoints[0])
}

func (this *EtcdUi) ScannerPort(ipAndPort string) bool {
	rs := false
	//tcpaddr,_ := net.ResolveTCPAddr("tcp4",ipAndPort)
	//_,err := net.DialTCP("tcp",nil,tcpaddr)
	_, err := net.DialTimeout("tcp", ipAndPort, 500*time.Millisecond)
	if err == nil {
		rs = true
	}
	return rs
}

//CRUD
func (this *EtcdUi) AddLease(key, value string, ttl int64) error {
	if this.ScannerPort(this.Endpoints[0]) {
		this.InitClientConn()
		defer this.Close()

		resp, err := this.ClientConn.Grant(context.TODO(), ttl)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = this.ClientConn.Put(ctx, key, value, clientv3.WithLease(resp.ID))
		cancel()
		if err != nil {
			switch err {
			case context.Canceled:
				fmt.Printf("ctx is canceled by another routine: %v\n", err)
			case context.DeadlineExceeded:
				fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
			case rpctypes.ErrEmptyKey:
				fmt.Printf("client-side error: %v\n", err)
			default:
				fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
			}
		}
		return err
	}
	return errors.New(fmt.Sprintf("%s unreachable", this.Endpoints[0]))
}

//CRUD
func (this *EtcdUi) Add(key, value string) error {
	if this.ScannerPort(this.Endpoints[0]) {
		this.InitClientConn()
		defer this.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := this.ClientConn.Put(ctx, key, value)
		cancel()
		if err != nil {
			switch err {
			case context.Canceled:
				fmt.Printf("ctx is canceled by another routine: %v\n", err)
			case context.DeadlineExceeded:
				fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
			case rpctypes.ErrEmptyKey:
				fmt.Printf("client-side error: %v\n", err)
			default:
				fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
			}
		}
		return err
	}
	return errors.New(fmt.Sprintf("%s unreachable", this.Endpoints[0]))
}

func (this *EtcdUi) Delete(key string) error {
	if this.ScannerPort(this.Endpoints[0]) {
		this.InitClientConn()
		defer this.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := this.ClientConn.Delete(ctx, key)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("%s unreachable", this.Endpoints[0]))
}

func (this *EtcdUi) DeleteAll(key string) error {
	if this.ScannerPort(this.Endpoints[0]) {
		this.InitClientConn()
		defer this.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := this.ClientConn.Delete(ctx, key, clientv3.WithPrefix())
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("%s unreachable", this.Endpoints[0]))
}

/*
Jtopo 总关系图谱
*/
//根据顶级机构和所有数据进行递归 得到树形结构的json字符串
//获取所有tree table最终数据
func (this *EtcdUi) GetTreeByMapJtopo() ([]map[string]interface{}, error) {
	this.GetTreeByString()
	// for k,v := range this.Tree {
	// 	fmt.Println(k,v)
	// }

	top := []string{}
	for _, v := range this.Tree {
		if v["parentOrg"] == "null" {
			top = append(top, v["name"])
		}
	}
	// fmt.Println("top",top)
	rs := this.GetMapJtopo(top, this.Tree)
	return rs, nil
}

//判断现有tree集合里面有key没有
func (this *EtcdUi) HasKeyByTreeToGet(key string) (bool, map[string]string) {
	rs := false
	if len(this.Tree) == 0 {
		return false, nil
	}
	for _, k := range this.Tree {
		if value, ok := k["name"]; ok {
			if value == key {
				rs = true
				return rs, k
			}
		}
	}
	return rs, nil
}

//根据map生成json数组
func (this *EtcdUi) GetMapJtopo(top []string, all []map[string]string) []map[string]interface{} {
	result := []map[string]interface{}{}
	//获取顶级项目以及子项目
	for _, y := range top {
		tmp := map[string]interface{}{}
		tmp["name"] = strings.Split(y, "/")[len(strings.Split(y, "/"))-1]
		//判断是否有子项目
		if this.HasChild(y, this.Tree) {
			tmp["nodes"] = this.GetMapJtopo(this.ForeignKeys(y, this.Tree), this.Tree)
		} else {
			if ok, data := this.HasKeyByTreeToGet(y); ok {
				tmp["ttl"] = data["ttl"]
				tmp["version"] = data["version"]
			}
			if tmp["value"] != nil {
				tmp["value"] = fmt.Sprintf("%s,%s", tmp["value"].(string), y)
			} else {
				tmp["value"] = y
			}
		}
		result = append(result, tmp)
	}
	return result
}
