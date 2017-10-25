package operator

import (
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

type Informer interface {
	//cache.SharedIndexInformer
	onAdd(interface{})
	onDelete(interface{})
	onUpdate(i,j interface{})
}

// HyperpilotOpertor watches the kubernetes api for changes to Pods and
// delete completed Pods without specific annotation
type HyperpilotOpertor struct {
	podInformer PodInformer
	deployInformer DeploymentInformer
	daemonSetInformer DaemonSetInformer
	kclient     *kubernetes.Clientset

	// map between resource type <-> controller
	//registeredController *sync.Map

	podRegisters 		[]IController
	deployRegisters 	[]IController
	daemonSetRegisters 	[]IController
	nsRegisters			[]IController
}


// NewHyperpilotOperator creates a new NewHyperpilotOperator
func NewHyperpilotOperator(kclient *kubernetes.Clientset, opts map[string]string) *HyperpilotOpertor {
	hpc := &HyperpilotOpertor{
		podRegisters: 		make([]IController, 0),
		deployRegisters: 	make([]IController, 0),
		daemonSetRegisters: make([]IController, 0),
		nsRegisters: 		make([]IController, 0),
		kclient:            kclient,
	}


	hpc.podInformer = InitPodInformer(kclient, opts, hpc)
	hpc.deployInformer = InitDeploymentInformer(kclient, opts, hpc)
	hpc.daemonSetInformer = InitDaemonSetInformer(kclient, opts, hpc)
	return hpc
}

// Run starts the process for listening for pod changes and acting upon those changes.
func (c *HyperpilotOpertor) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Printf("Listening for changes...")
	// When this function completes, mark the go function as done
	defer wg.Done()

	// Increment wait group as we're about to execute a go function
	wg.Add(1)

	// Execute go function
	go c.podInformer.indexInformer.Run(stopCh)
	go c.deployInformer.indexInformer.Run(stopCh)
	//go c.daemonSetInformer.indexInformer.Run(stopCh)

	// Wait till we receive a stop signal
	<-stopCh
}

//todo: slice is not thread-safe
func (c *HyperpilotOpertor)Accept(s IController, res resourceEnum) {

	if res.IsRegister(POD){
		c.podRegisters = append(c.podRegisters, s)
		log.Printf("Contoller {%s} registered resource POD", s)
	}

	if res.IsRegister(DEPLOYMENT){
		c.deployRegisters = append(c.deployRegisters, s)
		log.Printf("Contoller {%s} registered resource DEPLOYMENT", s)
	}

	if res.IsRegister(DAEMONSET){
		c.daemonSetRegisters = append(c.daemonSetRegisters, s)
		log.Printf("Contoller {%s} registered resource DAEMONSET", s)
	}

}

func (c *HyperpilotOpertor)Dcommission(id string) {
	//c.registeredController.Delete(id)
}