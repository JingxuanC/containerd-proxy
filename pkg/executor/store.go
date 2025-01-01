package executor

// The Store interface defines methods for managing pod and container information in a storage system.
// @property WritePodSandboxInfo - The `WritePodSandboxInfo` method in the `Store` interface is used to
// store information about a pod sandbox identified by its unique ID (`podUID`). The method takes the
// pod UID and a pointer to a `PodSandboxInfo` struct as parameters and stores this information in the
// underlying
// @property WriteContainerInfo - The `WriteContainerInfo` method in the `Store` interface is used to
// store information about a container identified by its `containerId`. The method takes the
// `containerId` as a string parameter and the corresponding `ContainerInfo` object as another
// parameter. This information can include details about the container
// @property GetPodSandboxInfo - The `GetPodSandboxInfo` method in the `Store` interface is used to
// retrieve information about a pod sandbox based on its unique identifier (`podUID`). It returns a
// pointer to a `PodSandboxInfo` struct that contains details about the pod sandbox.
// @property GetContainerInfo - The `GetContainerInfo` method in the `Store` interface is used to
// retrieve information about a specific container based on its unique identifier `containerUID`. This
// method should return a pointer to the `ContainerInfo` struct associated with the provided
// `containerUID`.
// @property DeletePodSandboxInfo - The `DeletePodSandboxInfo` method in the `Store` interface is used
// to remove the information associated with a specific pod sandbox identified by its unique identifier
// (`podUID`) from the store. This method allows for the cleanup of data related to a pod sandbox that
// is no longer needed or has
// @property DeleteContainerInfo - The `DeleteContainerInfo` method in the `Store` interface is used to
// remove the information associated with a specific container identified by its unique container ID
// from the store. This method allows for the deletion of container information stored in the data
// structure that implements the `Store` interface.
// @property ListPods - ListPods is a method in the Store interface that takes a function as an
// argument. This method is used to iterate over all the PodSandboxInfo objects stored in the data
// store and pass them to the provided function. The function receives a map where the keys are the
// unique identifiers of the pods
// @property ListContainers - The `ListContainers` method in the `Store` interface allows you to
// iterate over all the stored container information by providing a function that takes a map of
// container IDs to `ContainerInfo` objects as an argument. This method enables you to access and work
// with the container information stored in the `Store
// @property {[]string} GetPodsID - GetPodsID is a method of the Store interface that returns a slice
// of strings containing the IDs of all the pods stored in the store.
// @property {[]string} GetContainersID - GetContainersID is a method defined in the Store interface
// that returns a slice of strings containing the IDs of all containers stored in the data store.
type Store interface {
	WritePodSandboxInfo(podUID string, pod *PodSandboxInfo)
	WriteContainerInfo(containerId string, container *ContainerInfo)
	GetPodSandboxInfo(podUID string) *PodSandboxInfo
	GetContainerInfo(containerUID string) *ContainerInfo
	DeletePodSandboxInfo(podUID string)
	DeleteContainerInfo(containerUID string)
	ListPods(func(m map[string]*PodSandboxInfo))
	ListContainers(func(m map[string]*ContainerInfo))
	GetPodsID() []string
	GetContainersID() []string
}
