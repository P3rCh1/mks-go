/*
Package cluster provides the ability to retrieve and manage Kubernetes clusters
through the MKS V2 API.

Example of getting a single cluster referenced by its id

	mksCluster, err := cluster.Get(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster.JSON200)

Example of getting all clusters

	mksClusters, err := cluster.List(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	for _, mksCluster := range mksClusters.JSON200.Items {
		fmt.Printf("%+v\n", mksCluster)
	}

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
	fmt.Printf("%+v\n", mksCluster.JSON201)

Example of updating an existing cluster

	updateOpts := &mksclient.ClusterUpdateStruct{
		MaintenanceWindowStart: new("07:00:00"),
		KubernetesOptions: &mksclient.KubernetesOptions{
			EnablePodSecurityPolicy: new(false),
			FeatureGates: &[]string{
				"TTLAfterFinished",
			},
			AdmissionControllers: &[]string{
				"NamespaceLifecycle",
			},
		},
	}
	mksCluster, err := cluster.Update(ctx, client, clusterID, updateOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster.JSON200)

Example of deleting a single cluster

	resp, err := cluster.Delete(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}

Example of getting a kubeconfig referenced by cluster id

	kubeconfig, err := cluster.GetKubeconfig(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(kubeconfig.Body))

Example of getting fields from Kubeconfig referenced by cluster id

	parsedKubeconfig, err := cluster.GetParsedKubeconfig(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server IP:", parsedKubeconfig.Server)
	fmt.Println("Cluster CA:", parsedKubeconfig.ClusterCA)
	fmt.Println("Client cert:", parsedKubeconfig.ClientCert)
	fmt.Println("Client key:", parsedKubeconfig.ClientKey)
	fmt.Println("Raw kubeconfig:", parsedKubeconfig.KubeconfigRaw)

Example of rotating certificates by cluster id

	resp, err := cluster.RotateCerts(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}

Example of upgrading Kubernetes patch version by cluster id

	mksCluster, err := cluster.UpgradePatchVersion(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster.JSON200)

Example of upgrading Kubernetes minor version by cluster id

	mksCluster, err := cluster.UpgradeMinorVersion(ctx, client, clusterID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mksCluster.JSON200)
*/
package cluster
