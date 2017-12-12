package etcd

/**
	y.Key
	y.CreateRevision
	y.Lease
	y.ModRevision
	y.Value
	y.Version
*/
func GetEtcdTemplate() []map[string]interface{} {
	column := []map[string]interface{}{}
	column = append(column,
		//map[string]interface{}{"field":"state","checkbox":true,"align":"center","valign":"middle"},
		map[string]interface{}{"field":"id","title":"Key","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"value","title":"Value","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"version","title":"version","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"lease","title":"Lease","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"createrevision","title":"CreateRevision","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"moderevision","title":"ModRevision","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"op","title":"操作","sortable":"false","align":"center","valign":"middle"},
	)
	return column
}

func GetEtcdServiceTemplate() []map[string]interface{} {
	column := []map[string]interface{}{}
	column = append(column,
		// map[string]interface{}{"field":"state","checkbox":true,"align":"center","valign":"middle"},
		map[string]interface{}{"field":"mgm","title":"管理","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"did","title":"ID","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dimage","title":"IMAGE","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dcommand","title":"Command","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dstate","title":"State","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dStatus","title":"Status","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dports","title":"Ports","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"dname","title":"Names","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"id","title":"Key","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"status","title":"状态","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"value","title":"Value","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"version","title":"version","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"lease","title":"Lease","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"createrevision","title":"CreateRevision","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"moderevision","title":"ModRevision","sortable":"true","align":"center","valign":"middle"},
		map[string]interface{}{"field":"op","title":"操作","sortable":"false","align":"center","valign":"middle"},
	)
	return column
}