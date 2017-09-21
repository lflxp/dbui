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
	)
	return column
}