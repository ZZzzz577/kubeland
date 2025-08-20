package kube

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/controller-gen crd paths="./..." output:crd:artifacts:config=config/crd
//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/controller-gen object paths="./..."
