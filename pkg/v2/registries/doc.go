/*
Package registries provides the ability to retrieve and manage cluster registries
through the MKS V2 API.

Example of getting cluster registries

	clusterRegistries, err := registries.Get(ctx, mksClient, clusterID)
	if err != nil {
	  log.Fatal(err)
	}
	for _, registry := range clusterRegistries {
	  fmt.Printf("%+v\n", registry)
	}
*/
package registries
