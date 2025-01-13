package k8s

import (
	"context"
	"strconv"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// creates two netpols
// 1. default-deny-all: block all ingress and egress traffic for all pods in the configured namespace
// 2. default-allow-namespaces: allow all ingress and egress traffic to pods in the same namespace
func CreateDefaultNetpols(namespace string) error {

	// create a network policy block all ingress and egress traffic
	// for all pods in the configured namespace
	client, err := getKubeClient()
	if err != nil {
		return err
	}

	// create a network policy that allows all ingress and egress traffic from the same namespace
	allowNsInternal := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-allow-same-namespace",
			Namespace: namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress, networkingv1.PolicyTypeEgress},
			PodSelector: metav1.LabelSelector{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"kubernetes.io/metadata.name": namespace,
								},
							},
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"kubernetes.io/metadata.name": namespace,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = client.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), allowNsInternal, metav1.CreateOptions{})

	return nil

}

func CreateEgressNetpol(namespace string, egressEndpoints []string) error {
	if len(egressEndpoints) == 0 || egressEndpoints[0] == "" {
		return nil
	}

	// create a network policy that allows egress traffic to the specified endpoints
	client, err := getKubeClient()
	if err != nil {
		return err
	}

	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-egress",
			Namespace: namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeEgress},
			PodSelector: metav1.LabelSelector{},
			Egress:      []networkingv1.NetworkPolicyEgressRule{},
		},
	}

	for _, egress := range egressEndpoints {

		e := strings.Split(egress, ":")

		port, err := strconv.Atoi(e[1])
		if err != nil {
			return err
		}

		np.Spec.Egress = append(np.Spec.Egress, networkingv1.NetworkPolicyEgressRule{
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(port),
					},
				},
			},
			To: []networkingv1.NetworkPolicyPeer{
				{
					IPBlock: &networkingv1.IPBlock{
						CIDR: e[0],
					},
				},
			},
		})
	}

	_, err = client.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), np, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
