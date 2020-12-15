package resource

import (
	"sigs.k8s.io/yaml"
)

// Resource represents a Kubernetes resource within a file
type Resource struct {
	Path  string
	Bytes []byte
	sig   *Signature
}

// Signature is a key representing a Kubernetes resource
type Signature struct {
	Kind, Version, Namespace, Name string
}

// Signature computes a signature for a resource, based on its Kind, Version, Namespace & Name
func (res *Resource) Signature() (*Signature, error) {
	if res.sig != nil {
		return res.sig, nil
	}

	resource := struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name         string `yaml:"name"`
			Namespace    string `yaml:"namespace"`
			GenerateName string `yaml:"generateName"`
		} `yaml:"Metadata"`
	}{}
	err := yaml.Unmarshal(res.Bytes, &resource)

	name := resource.Metadata.Name
	if resource.Metadata.GenerateName != "" {
		name = resource.Metadata.GenerateName + "{{ generateName }}"
	}

	// We cache the result to not unmarshall every time we want to access the signature
	res.sig = &Signature{Kind: resource.Kind, Version: resource.APIVersion, Namespace: resource.Metadata.Namespace, Name: name}
	return res.sig, err
}

func (res *Resource) SignatureFromMap(m map[string]interface{}) (*Signature, error) {
	if res.sig != nil {
		return res.sig, nil
	}

	APIVersion, _ := m["apiVersion"].(string)
	Kind, _ := m["kind"].(string)

	var name, ns string
	Metadata, ok := m["metadata"].(map[string]interface{})
	if ok {
		name, _ = Metadata["name"].(string)
		ns, _ = Metadata["namespace"].(string)
		if _, ok := Metadata["generateName"].(string); ok {
			name = Metadata["generateName"].(string) + "{{ generateName }}"
		}
	}

	// We cache the result to not unmarshall every time we want to access the signature
	res.sig = &Signature{Kind: Kind, Version: APIVersion, Namespace: ns, Name: name}
	return res.sig, nil
}
