package executor

import (
	"sync"

	"github.com/dolthub/swiss"
)

type MemoryStore struct {
	plock sync.RWMutex
	pods  *swiss.Map[string, *PodSandboxInfo]

	clock      sync.RWMutex
	containers *swiss.Map[string, *ContainerInfo]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pods:       swiss.NewMap[string, *PodSandboxInfo](2048),
		containers: swiss.NewMap[string, *ContainerInfo](2048),
	}
}

func (m *MemoryStore) WritePodSandboxInfo(podUID string, pod *PodSandboxInfo) {
	m.plock.Lock()
	defer m.plock.Unlock()
	m.pods.Put(podUID, pod)
}

func (m *MemoryStore) WriteContainerInfo(containerId string, container *ContainerInfo) {
	m.clock.Lock()
	defer m.clock.Unlock()
	m.containers.Put(containerId, container)
}

func (m *MemoryStore) GetPodSandboxInfo(podUID string) *PodSandboxInfo {
	m.plock.RLock()
	defer m.plock.RUnlock()
	pod, ok := m.pods.Get(podUID)
	if !ok {
		return nil
	}
	return pod
}

func (m *MemoryStore) GetContainerInfo(containerUID string) *ContainerInfo {
	m.clock.RLock()
	defer m.clock.RUnlock()
	container, ok := m.containers.Get(containerUID)
	if !ok {
		return nil
	}
	return container
}

func (m *MemoryStore) DeletePodSandboxInfo(podUID string) {
	m.plock.Lock()
	defer m.plock.Unlock()
	m.pods.Delete(podUID)
}

func (m *MemoryStore) DeleteContainerInfo(containerUID string) {
	m.clock.Lock()
	defer m.clock.Unlock()
	m.containers.Delete(containerUID)
}

func (m *MemoryStore) ListPods(f func(m map[string]*PodSandboxInfo)) {
	m.plock.RLock()
	defer m.plock.RUnlock()

	pods := make(map[string]*PodSandboxInfo)
	m.pods.Iter(func(k string, v *PodSandboxInfo) (stop bool) {
		pods[k] = v
		return false
	})

	f(pods)

}

func (m *MemoryStore) ListContainers(f func(m map[string]*ContainerInfo)) {
	m.clock.RLock()
	defer m.clock.RUnlock()

	containers := make(map[string]*ContainerInfo)
	m.containers.Iter(func(k string, v *ContainerInfo) (stop bool) {
		containers[k] = v
		return false
	})

	f(containers)
}

func (m *MemoryStore) GetPodsID() []string {
	m.plock.RLock()
	defer m.plock.RUnlock()

	pods := make([]string, 0, m.pods.Count())
	m.pods.Iter(func(k string, v *PodSandboxInfo) (stop bool) {
		pods = append(pods, k)
		return false
	})

	return pods
}

func (m *MemoryStore) GetContainersID() []string {
	m.clock.RLock()
	defer m.clock.RUnlock()

	containers := make([]string, 0, m.containers.Count())
	m.containers.Iter(func(k string, v *ContainerInfo) (stop bool) {
		containers = append(containers, k)
		return false
	})

	return containers
}
