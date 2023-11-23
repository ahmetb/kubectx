package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) clustersNode() (*yaml.Node, error) {
	clusters := valueOf(k.rootNode, "clusters")
	if clusters == nil {
		return nil, errors.New("\"clusters\" entry is nil")
	} else if clusters.Kind != yaml.SequenceNode {
		return nil, errors.New("\"clusters\" is not a sequence node")
	}
	return clusters, nil
}

func (k *Kubeconfig) ClusterOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}

	return k.clusterOfContextNode(ctx)
}

func (k *Kubeconfig) clusterOfContextNode(contextNode *yaml.Node) (string, error) {
	ctxBody := valueOf(contextNode, "context")
	if ctxBody == nil {
		return "", errors.New("no context field found for context entry")
	}

	cluster := valueOf(ctxBody, "cluster")
	if cluster == nil || cluster.Value == "" {
		return "", errors.New("no cluster field found for context entry")
	}
	return cluster.Value, nil
}

func (k *Kubeconfig) CountClusterReferences(clusterName string) (int, error) {
	contexts, err := k.contextsNode()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, contextNode := range contexts.Content {
		contextCluster, err := k.clusterOfContextNode(contextNode)
		if err != nil {
			return 0, err
		}
		if clusterName == contextCluster {
			count += 1
		}
	}

	return count, nil
}

func (k *Kubeconfig) DeleteClusterEntry(deleteName string) error {
	contexts, err := k.clustersNode()
	if err != nil {
		return err
	}

	deleteNamedChildNode(contexts, deleteName)
	return nil
}
