package executor

import (
	"fmt"
	"reflect"

	v1 "github.com/JingxuanC/containerd-proxy/apis/v1"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
)

type PodSandboxInfo struct {
	*v1.PodSandboxHookRequest
}

type PodExecutor struct {
	store Store
	*PodSandboxInfo
}

func NewPodExecutor(s Store) *PodExecutor {
	return &PodExecutor{
		store: s,
		PodSandboxInfo: &PodSandboxInfo{
			PodSandboxHookRequest: &v1.PodSandboxHookRequest{},
		},
	}
}

func (p *PodExecutor) String() string {
	return fmt.Sprintf("pod(%v/%v)", p.GetPodMeta().GetName(), p.GetPodMeta().Uid)
}

func (p *PodExecutor) GetMetaInfo() string {
	return fmt.Sprintf("pod(%v/%v)", p.GetPodMeta().GetName(), p.GetPodMeta().Uid)
}

func (p *PodExecutor) GenerateHookRequest() interface{} {
	return p.PodSandboxHookRequest
}

func (p *PodExecutor) ParseRequest(req interface{}) (string, error) {
	switch request := req.(type) {
	case *runtimeapi.RunPodSandboxRequest:
		p.PodSandboxHookRequest = &v1.PodSandboxHookRequest{
			PodMeta: &v1.PodSandboxMetadata{
				Name:    request.GetConfig().GetMetadata().GetName(),
				Uid:     request.GetConfig().GetMetadata().GetUid(),
				Attempt: request.GetConfig().GetMetadata().GetAttempt(),
			},
			RuntimeHandler: request.GetRuntimeHandler(),
			Annotations:    request.GetConfig().GetAnnotations(),
			Labels:         request.GetConfig().GetLabels(),
			CgroupParent:   criCgroupParent(request.GetConfig()),
			Resources:      CRIResource(request.GetConfig()),
		}
	case *runtimeapi.StopPodSandboxRequest:
		return request.GetPodSandboxId(), p.GetPodSandbox(request.GetPodSandboxId())
	}
	return "", nil
}

func (P *PodExecutor) GetPodSandbox(podId string) error {
	pod := P.store.GetPodSandboxInfo(podId)
	if pod == nil {
		return fmt.Errorf("pod %v not found", podId)
	}
	P.PodSandboxInfo = pod
	return nil
}

func (p *PodExecutor) ResourceCheckPoint(response interface{}) error {
	resp, ok := response.(*runtimeapi.RunPodSandboxResponse)
	if !ok || p.PodSandboxHookRequest == nil {
		return nil
	}

	klog.Infof(" store pod(%v/%v) checkpoint", p.GetPodMeta().GetName(), p.GetPodMeta().Uid)
	p.store.WritePodSandboxInfo(resp.GetPodSandboxId(), p.PodSandboxInfo)
	return nil
}

func (p *PodExecutor) DeleteCheckpointIfNeed(request interface{}) error {
	switch request.(type) {
	case *runtimeapi.StopPodSandboxRequest:
		klog.Infof(" delete pod(%v/%v) checkpoint", p.GetPodMeta().GetName(), p.GetPodMeta().Uid)
		p.store.DeletePodSandboxInfo(p.GetPodMeta().GetUid())
	}
	return nil
}
func (p *PodExecutor) DeleteCheckpointForce(id string) {
	klog.Infof(" delete pod(%v/%v) checkpoint", p.GetPodMeta().GetName(), p.GetPodMeta().Uid)
	p.store.DeletePodSandboxInfo(id)
}
func (p *PodExecutor) UpdateRequest(response interface{}, request interface{}) error {
	resp, ok := response.(*v1.PodSandboxHookResponse)

	if !ok {
		return fmt.Errorf("response is not *runtimeapi.RunPodSandboxResponse, but get %s", reflect.TypeOf(response).String())
	}
	// update PodRecourcesExecutor
	p.Annotations = MergeLables(p.Annotations, resp.Annotations)
	p.Labels = MergeLables(p.Labels, resp.Labels)

	if resp.CgroupParent != "" {
		p.CgroupParent = resp.CgroupParent
	}

	//update cri request
	switch req := request.(type) {
	case *runtimeapi.RunPodSandboxRequest:
		if p.Annotations != nil {
			req.Config.Annotations = p.Annotations
		}
		if p.Labels != nil {
			req.Config.Labels = p.Labels
		}
		if p.CgroupParent != "" {
			req.Config.Linux.CgroupParent = p.CgroupParent
		}
		if p.Resources != nil {
			req.Config.Linux.Resources = ToCRIContainerResources(p.Resources)
		}

	}

	return nil

}
