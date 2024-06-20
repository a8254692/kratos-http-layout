package registryx

//// NewKubeRegistry ...
//func NewKubeRegistry() (registry.Registrar, error) {
//	config, err := rest.InClusterConfig()
//	if err != nil {
//		return nil, err
//	}
//	clientSet, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return nil, err
//	}
//	reg := kuberegistry.NewRegistry(clientSet)
//	reg.Start()
//	return reg, nil
//}
//
//// NewKubeDiscovery ...
//func NewKubeDiscovery() (registry.Discovery, error) {
//	config, err := rest.InClusterConfig()
//	if err != nil {
//		return nil, err
//	}
//	clientSet, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return nil, err
//	}
//	discovery := kuberegistry.NewRegistry(clientSet)
//	discovery.Start()
//	return discovery, nil
//}
