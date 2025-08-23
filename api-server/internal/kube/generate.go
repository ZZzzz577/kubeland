package kube

//go:generate bash ./hack/codegen.sh
//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/controller-gen crd paths="./..." output:crd:artifacts:config=config/crd
