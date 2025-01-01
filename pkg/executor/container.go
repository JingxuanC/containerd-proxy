package executor

import (
	"fmt"

	v1 "github.com/JingxuanC/containerd-proxy/apis/v1"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// ContainerInfo is almost the same with v1alpha.ContainerResourceHookRequest
type ContainerInfo struct {
	*v1.ContainerResourceHookRequest
}

type ContainerExecuotor struct {
	store Store
	ContainerInfo
}

func NewContainerExecutor(s Store) *ContainerExecuotor {
	return &ContainerExecuotor{
		store: s,
		ContainerInfo: ContainerInfo{
			ContainerResourceHookRequest: &v1.ContainerResourceHookRequest{},
		},
	}
}

func (c *ContainerExecuotor) String() string {
	return fmt.Sprintf("pod(%v/%v)container(%v)",
		c.GetPodMeta().GetName(), c.GetPodMeta().Uid, c.GetContainerMeta().GetName())
}

func (c *ContainerExecuotor) GetMetaInfo() string {
	return fmt.Sprintf("pod(%v/%v)container(%v)",
		c.GetPodMeta().GetName(), c.GetPodMeta().GetUid(),
		c.GetContainerMeta().GetName())
}

func (c *ContainerExecuotor) GenerateHookRequest() interface{} {
	return c.ContainerResourceHookRequest
}

// The `ParseRequest` method in the `ContainerExecuotor` struct is used to parse different types of
// runtime requests and extract relevant information from them.
func (c *ContainerExecuotor) ParseRequest(req interface{}) (string, error) {
	switch request := req.(type) {
	case *runtimeapi.CreateContainerRequest:
		podID := request.GetPodSandboxId()
		podCheckPoint := c.store.GetPodSandboxInfo(podID)
		if podCheckPoint == nil {
			return "", fmt.Errorf("pod(%v) checkpoint not found", podID)
		}
		c.ContainerInfo = ContainerInfo{
			ContainerResourceHookRequest: &v1.ContainerResourceHookRequest{
				PodMeta:              podCheckPoint.PodMeta,
				PodResources:         podCheckPoint.Resources,
				PodAnnotations:       podCheckPoint.Annotations,
				PodLabels:            podCheckPoint.Labels,
				ContainerAnnotations: request.GetConfig().GetAnnotations(),
				ContainerMeta: &v1.ContainerMetadata{
					Name:    request.GetConfig().GetMetadata().GetName(),
					Attempt: request.GetConfig().GetMetadata().GetAttempt(),
				},
				ContainerEnvs:      ToAPIsContainerEnvs(request.GetConfig().GetEnvs()),
				PodCgroupParent:    podCheckPoint.CgroupParent,
				ContainerResources: ToAPIsContainerResource(request.GetConfig().GetLinux().GetResources()),
			},
		}
		return podID, nil

	case *runtimeapi.StartContainerRequest:
		return "", nil
	}
	return "", nil
}

func (c *ContainerExecuotor) LoadContainerInfo(containerId, stage string) error {
	containerCheckPoint := c.store.GetContainerInfo(containerId)

	if containerCheckPoint == nil {
		return fmt.Errorf("fail to load contaier(%v) from store during %v", containerId, stage)
	}
	c.ContainerInfo = *containerCheckPoint

	return nil
}

func (c *ContainerExecuotor) DeleteCheckpointIfNeed(req interface{}) error {
	switch request := req.(type) {
	case *runtimeapi.StopContainerRequest:
		c.store.DeleteContainerInfo(request.GetContainerId())
		return nil
	case *runtimeapi.RemoveContainerRequest:
		c.store.DeleteContainerInfo(request.GetContainerId())
		return nil
	default:
		return fmt.Errorf("unsupported request type %T", request)
	}
}
