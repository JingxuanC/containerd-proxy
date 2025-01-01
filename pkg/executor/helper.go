package executor

import (
	v1 "github.com/JingxuanC/containerd-proxy/apis/v1"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func ToAPIsContainerResource(r *runtimeapi.LinuxContainerResources) *v1.LinuxContainerResources {

	linuxResouces := &v1.LinuxContainerResources{
		CpuPeriod:              r.GetCpuPeriod(),
		CpuQuota:               r.GetCpuQuota(),
		CpuShares:              r.GetCpuShares(),
		CpusetCpus:             r.GetCpusetCpus(),
		CpusetMems:             r.GetCpusetMems(),
		MemoryLimitInBytes:     r.GetMemoryLimitInBytes(),
		MemorySwapLimitInBytes: r.GetMemorySwapLimitInBytes(),
		OomScoreAdj:            r.GetOomScoreAdj(),
		Unified:                r.GetUnified(),
	}

	for _, h := range r.GetHugepageLimits() {
		linuxResouces.HugepageLimits = append(linuxResouces.HugepageLimits, &v1.HugepageLimit{
			PageSize: h.GetPageSize(),
			Limit:    h.GetLimit(),
		})
	}
	return linuxResouces
}

func ToAPIsContainerEnvs(envs []*runtimeapi.KeyValue) map[string]string {
	envMap := make(map[string]string)
	if envs == nil {
		return envMap
	}
	for _, env := range envs {
		envMap[env.GetKey()] = env.GetValue()
	}
	return envMap
}

func ToCRIContainerResources(r *v1.LinuxContainerResources) *runtimeapi.LinuxContainerResources {

	linuxResouces := &runtimeapi.LinuxContainerResources{
		CpuPeriod:              r.GetCpuPeriod(),
		CpuQuota:               r.GetCpuQuota(),
		CpuShares:              r.GetCpuShares(),
		CpusetCpus:             r.GetCpusetCpus(),
		MemoryLimitInBytes:     r.GetMemoryLimitInBytes(),
		MemorySwapLimitInBytes: r.GetMemorySwapLimitInBytes(),
		OomScoreAdj:            r.GetOomScoreAdj(),
		Unified:                r.GetUnified(),
		CpusetMems:             r.GetCpusetMems(),
	}

	for _, h := range r.GetHugepageLimits() {
		linuxResouces.HugepageLimits = append(linuxResouces.HugepageLimits, &runtimeapi.HugepageLimit{
			PageSize: h.GetPageSize(),
			Limit:    h.GetLimit(),
		})
	}
	return linuxResouces
}

func UpdateAPIsContainerResource(r, resources *v1.LinuxContainerResources) *v1.LinuxContainerResources {
	if r == nil || resources == nil {
		return resources
	}

	if resources.CpuPeriod > 0 {
		r.CpuPeriod = resources.CpuPeriod
	}
	if resources.CpuQuota > 0 {
		r.CpuQuota = resources.CpuQuota
	}

	if resources.CpuShares > 0 {
		r.CpuShares = resources.CpuShares
	}
	if resources.MemoryLimitInBytes > 0 {
		r.MemoryLimitInBytes = resources.MemoryLimitInBytes
	}
	if resources.MemorySwapLimitInBytes > 0 {
		r.MemorySwapLimitInBytes = resources.MemorySwapLimitInBytes
	}
	if resources.OomScoreAdj >= -1000 && resources.OomScoreAdj <= 1000 {
		r.OomScoreAdj = resources.OomScoreAdj
	}

	r.CpusetCpus = resources.CpusetCpus
	r.CpusetMems = resources.CpusetMems
	r.Unified = resources.Unified

	return r
}

func CRIResource(config *runtimeapi.PodSandboxConfig) *v1.LinuxContainerResources {
	if config == nil {
		return nil
	}
	if config.Linux == nil {
		return nil
	}

	if config.Linux.Resources == nil {
		return nil
	}

	hlimits := make([]*v1.HugepageLimit, 0, 4)

	for _, i := range config.Linux.Resources.HugepageLimits {
		hlimits = append(hlimits, &v1.HugepageLimit{
			PageSize: i.GetPageSize(),
			Limit:    i.GetLimit(),
		})
	}

	return &v1.LinuxContainerResources{
		CpuPeriod:              config.Linux.Resources.CpuPeriod,
		CpuQuota:               config.Linux.Resources.CpuQuota,
		CpuShares:              config.Linux.Resources.CpuShares,
		CpusetCpus:             config.Linux.Resources.CpusetCpus,
		MemoryLimitInBytes:     config.Linux.Resources.MemoryLimitInBytes,
		MemorySwapLimitInBytes: config.Linux.Resources.MemorySwapLimitInBytes,
		OomScoreAdj:            config.Linux.Resources.OomScoreAdj,
		HugepageLimits:         hlimits,
		Unified:                config.Linux.Resources.Unified,
		CpusetMems:             config.Linux.Resources.CpusetMems,
	}

}

func criCgroupParent(config *runtimeapi.PodSandboxConfig) string {
	if config == nil {
		return ""
	}

	if config.Linux == nil {
		return ""
	}

	return config.Linux.CgroupParent
}
