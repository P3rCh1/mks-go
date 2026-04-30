/*
Package cluster provides the ability to retrieve and manage Kubernetes clusters
through the MKS V2 API.

Example of getting a single cluster referenced by its id

	mksCluster, err := cluster.Get(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster)

Example of creating a new cluster

	clusterBody := &mksclient.ClusterCreateStruct{
		Name:        "test-cluster",
		KubeVersion: "1.28",
		Region:      new("ru-1"),
		Nodegroups:  []*mksclient.NodegroupCreateStruct{
			{
				Name:             "test-nodegroup",
				Count:            1,
				CPUs:             new(2),
				RAM:              new(4096),
				VolumeSize:       new(50),
				VolumeType:       new("standard"),
				AvailabilityZone: new("ru-3a"),
			},
		},
	}
	mksCluster, err := cluster.Create(ctx, client, clusterBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster)
*/
package cluster
