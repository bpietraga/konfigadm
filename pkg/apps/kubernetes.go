package apps

import . "github.com/moshloop/konfigadm/pkg/types"

var Kubernetes Phase = kubernetes{}

type kubernetes struct{}

func (k kubernetes) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	if sys.Kubernetes == nil {
		return []Command{}, Filesystem{}, nil
	}

	sys.
		AppendPackageRepo(PackageRepo{
			Name:            "kubernetes",
			URL:             "https://apt.kubernetes.io/",
			VersionCodeName: "kubernetes-xenial",
			GPGKey:          "https://packages.cloud.google.com/apt/doc/apt-key.gpg",
		}, DEBIAN_LIKE).
		AppendPackageRepo(PackageRepo{
			Name:   "kubernetes",
			URL:    "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
			GPGKey: "https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
		}, REDHAT_LIKE).
		AppendPackageRepo(PackageRepo{
			Name:   "kubernetes",
			URL:    "https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64",
			GPGKey: "https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg",
		}, FEDORA)

	sys.AddPackage("kubelet-"+sys.Kubernetes.Version+"-0", &REDHAT_LIKE).
		AddPackage("kubeadm-"+sys.Kubernetes.Version+"-0", &REDHAT_LIKE).
		AddPackage("kubectl-"+sys.Kubernetes.Version+"-0", &REDHAT_LIKE)

	sys.AddPackage("kubelet-"+sys.Kubernetes.Version+"-0", &FEDORA).
		AddPackage("kubeadm-"+sys.Kubernetes.Version+"-0", &FEDORA).
		AddPackage("kubectl-"+sys.Kubernetes.Version+"-0", &FEDORA)

	sys.AddPackage("kubelet=="+sys.Kubernetes.Version+"-00", &DEBIAN_LIKE).
		AddPackage("kubeadm=="+sys.Kubernetes.Version+"-00", &DEBIAN_LIKE).
		AddPackage("kubectl=="+sys.Kubernetes.Version+"-00", &DEBIAN_LIKE)

	sys.AddPackage("socat jq ebtables ntp libseccomp nfs-utils", &REDHAT_LIKE)
	sys.AddPackage("socat jq ebtables ntp libseccomp2 nfs-client", &DEBIAN_LIKE)

	sys.Environment["KUBECONFIG"] = "/etc/kubernetes/admin.conf"
	sys.Sysctls["vm.swappiness"] = "0"
	sys.Sysctls["net.bridge.bridge-nf-call-iptables"] = "1"
	sys.Sysctls["net.bridge.bridge-nf-call-ip6tables"] = "1"
	sys.Sysctls["net.ipv4.ip_forward"] = "1"

	fs := Filesystem{}
	fs["/etc/modules-load.d/kubernetes.conf"] = File{Content: "overlay\nbr_netfilter"}

	//TODO(moshloop) ensure kernel modules are loaded: overlay and br_netfilter

	return []Command{}, fs, nil

}
