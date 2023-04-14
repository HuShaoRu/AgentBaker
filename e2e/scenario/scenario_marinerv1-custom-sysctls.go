package scenario

import (
	"fmt"
	"strings"

	"github.com/Azure/agentbaker/pkg/agent/datamodel"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

func marinerv1CustomSysctls() *Scenario {
	customNfConntrackMax := 300000
	customNfConntrackBuckets := 120320
	return &Scenario{
		Name:        "marinerv1-custom-sysctls",
		Description: "tests that a MarinerV1 VHD can be properly bootstrapped when supplied custom node config that contains custom sysctl settings",
		Config: Config{
			ClusterSelector: NetworkPluginKubenetSelector,
			ClusterMutator:  NetworkPluginKubenetMutator,
			BootstrapConfigMutator: func(nbc *datamodel.NodeBootstrappingConfiguration) {
				customLinuxConfig := &datamodel.CustomLinuxOSConfig{
					Sysctls: &datamodel.SysctlConfig{
						NetNetfilterNfConntrackMax:     to.Ptr(int32(customNfConntrackMax)),
						NetNetfilterNfConntrackBuckets: to.Ptr(int32(customNfConntrackBuckets)),
					},
				}
				nbc.AgentPoolProfile.CustomLinuxOSConfig = customLinuxConfig
				nbc.ContainerService.Properties.AgentPoolProfiles[0].CustomLinuxOSConfig = customLinuxConfig
				nbc.ContainerService.Properties.AgentPoolProfiles[0].Distro = "aks-cblmariner-v1"
				nbc.AgentPoolProfile.Distro = "aks-cblmariner-v1"
			},
			VMConfigMutator: func(vmss *armcompute.VirtualMachineScaleSet) {
				vmss.Properties.VirtualMachineProfile.StorageProfile.ImageReference = &armcompute.ImageReference{
					ID: to.Ptr(DefaultImageVersionIDs["marinerv1"]),
				}
			},
			LiveVMValidators: []*LiveVMValidator{
				{
					Description: "assert custom sysctl settings",
					Command:     `sysctl -a | grep "net.netfilter"`,
					Asserter: func(stdout, stderr string) error {
						if !strings.Contains(stdout, fmt.Sprintf("net.netfilter.nf_conntrack_max = %d", customNfConntrackMax)) {
							return fmt.Errorf(fmt.Sprintf("expected to find net.netfilter.nf_conntrack_max set to %d, but was not", customNfConntrackMax))
						}
						if !strings.Contains(stdout, fmt.Sprintf("net.netfilter.nf_conntrack_buckets = %d", customNfConntrackBuckets)) {
							return fmt.Errorf(fmt.Sprintf("expected to find net.netfilter.nf_conntrack_buckets set to %d, but was not", customNfConntrackBuckets))
						}
						return nil
					},
				},
			},
		},
	}
}
