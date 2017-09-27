package qhandler_elasticsearch

type Settings struct {
	NumShards int "json:`index.number_of_shards`"
	NumReplicas int "json:`index.number_of_replicas`"
}

/*type Mappings struct {
	NumShards int "json:`index.number_of_shards`"
	NumReplicas int "json:`index.number_of_replicas`"
}*/
