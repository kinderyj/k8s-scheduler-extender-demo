package algorithm

import (
	"math/rand"
	"strings"

	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

const (
	// LuckyPred rejects a node if you're not lucky ¯\_(ツ)_/¯
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Well, you're not lucky"
	ISPPred          = "ISPPred"
)

var predicatesFuncs = map[string]FitPredicate{
	LuckyPred: LuckyPredicate,
	ISPPred:   ISPPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{ISPPred}

// filter filters nodes according to predicates defined in this extender
// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFitPod()
func Filter(args extender.ExtenderArgs) *extender.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(extender.FailedNodesMap)
	pod := args.Pod

	// TODO: parallelize this
	// TODO: handle error
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := extender.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	var failReasons []string
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := rand.Intn(2) == 0
	if lucky {
		klog.Infof("pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	klog.Infof("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{LuckyPredFailMsg}, nil
}

func ISPPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	if isp, ok := pod.Labels["isp"]; ok {
		klog.Infof("pod %v/%v needs the isp %v, try to check the isp the node %v", pod.Name, pod.Namespace, isp, node.Name)
		if nodeISP, ok := node.Labels["isp"]; ok {
			if isp == nodeISP {
				klog.Infof("pod %v/%v fits on node %v according to isp predicate\n", pod.Name, pod.Namespace, node.Name)
				return true, nil, nil
			}
			klog.Infof("pod %v/%v doesn't fits on node %v according to isp predicate\n", pod.Name, pod.Namespace, node.Name)
			return false, []string{"Nodes's isp lable doesn't fit pods's isp label."}, nil
		}
		klog.Infof("pod %v/%v doesn't fits on node %v according to isp predicate\n", pod.Name, pod.Namespace, node.Name)
		return false, []string{"Nodes doesn't have any isp labels."}, nil

	}
	klog.Infof("pod %v/%v doesn't specified the isp, ingnore the ISP Predicate.", pod.Name, pod.Namespace)
	return true, nil, nil
}
