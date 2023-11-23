package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) usersNode() (*yaml.Node, error) {
	users := valueOf(k.rootNode, "users")
	if users == nil {
		return nil, errors.New("\"users\" entry is nil")
	} else if users.Kind != yaml.SequenceNode {
		return nil, errors.New("\"users\" is not a sequence node")
	}
	return users, nil
}

func (k *Kubeconfig) UserOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}

	return k.userOfContextNode(ctx)
}

func (k *Kubeconfig) userOfContextNode(contextNode *yaml.Node) (string, error) {
	ctxBody := valueOf(contextNode, "context")
	if ctxBody == nil {
		return "", errors.New("no context field found for context entry")
	}

	user := valueOf(ctxBody, "user")
	if user == nil || user.Value == "" {
		return "", errors.New("no user field found for context entry")
	}
	return user.Value, nil
}

func (k *Kubeconfig) CountUserReferences(userName string) (int, error) {
	contexts, err := k.contextsNode()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, contextNode := range contexts.Content {
		contextUser, err := k.userOfContextNode(contextNode)
		if err != nil {
			return 0, err
		}
		if userName == contextUser {
			count += 1
		}
	}

	return count, nil
}

func (k *Kubeconfig) DeleteUserEntry(deleteName string) error {
	contexts, err := k.usersNode()
	if err != nil {
		return err
	}

	deleteNamedChildNode(contexts, deleteName)
	return nil
}
