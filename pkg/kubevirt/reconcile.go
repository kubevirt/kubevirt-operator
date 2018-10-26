package kubevirt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/kubevirt/kubevirt-operator/pkg/apis/virt/v1alpha1"

	"github.com/sirupsen/logrus"
)

type config struct {
	version  string
	registry string
	manifest string
}

// Reconcile - Reconcile the current state of the world with the expected state
func Reconcile(v *v1alpha1.Virt, deleted bool) error {
	conf, err := newConfig(v)
	if err != nil {
		return err
	}

	err = renderConfigFile(conf)
	if err != nil {
		return err
	}

	if deleted {
		consumeManifest(conf, "delete")
		return nil
	}
	consumeManifest(conf, "apply")

	return nil
}

func consumeManifest(c config, action string) {
	// Create kubevirt manifest using the client
	cmd := exec.Command("kubectl", action, "-f", c.manifest)
	// Error is outputed in plain text in out
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "Error from server (AlreadyExists)") {
		logrus.Infof("KubeVirt resources already exist in the cluster")
		logrus.Debugf(string(out))
	} else {
		logrus.Errorf("Applying KubeVirt %s manifest", c.version)
		logrus.Infof(string(out))
	}
}

// Until we have go templates for manifests, use string replace
func renderConfigFile(c config) error {
	in, err := ioutil.ReadFile(c.manifest)
	if err != nil {
		logrus.Errorf("KubeVirt Version %s is not supported by the operator", c.version)
		return err
	}

	lines := strings.Split(string(in), "\n")
	for l, line := range lines {
		if strings.Contains(line, "docker.io") {
			// Replace registry
			r := strings.Replace(line, "docker.io", c.registry, -1)
			lines[l] = r
		}
	}
	o := strings.Join(lines, "\n")
	err = ioutil.WriteFile(c.manifest, []byte(o), 0644)
	if err != nil {
		logrus.Errorf("Failed to render new config file")
		return err
	}

	return nil
}

func newConfig(v *v1alpha1.Virt) (config, error) {

	kubevirtVersion := v.Spec.Version
	registry := v.Spec.Registry
	if kubevirtVersion == "" {
		kubevirtVersion = LatestKubevirtVersion
	}
	if registry == "" {
		registry = "docker.io"
	}

	manifest := fmt.Sprintf("/etc/kubevirt/%s/kubevirt.yaml", kubevirtVersion)

	_, err := os.Stat(manifest)
	if os.IsNotExist(err) {
		logrus.Errorf("KubeVirt Version %s is not supported by the operator", kubevirtVersion)
	}
	if err != nil {
		return config{}, errors.New("Unsupported KubeVirt Version")
	}

	return config{
		version:  kubevirtVersion,
		registry: registry,
		manifest: manifest,
	}, nil
}
