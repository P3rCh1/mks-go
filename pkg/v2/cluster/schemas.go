package cluster

// Kubeconfig represents the structure of parsed kubeconfig file.
type Kubeconfig struct {
	Clusters []struct {
		Cluster struct {
			Server                   string `yaml:"server"`
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

// KubeconfigInfo represents a parsed kubeconfig file.
type KubeconfigInfo struct {
	ClusterCA     string
	Server        string
	ClientCert    string
	ClientKey     string
	KubeconfigRaw string
}
