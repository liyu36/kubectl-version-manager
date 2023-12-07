package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func init() {
	if !KubectlIsExist(GetBaseDir()) {
		if err := os.MkdirAll(GetBaseDir(), os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
}

func JoinSchema(s string) string {
	return fmt.Sprintf("https://%s", s)
}

const kubernetesSourceGitURL = "github.com/kubernetes/kubernetes"

func GetkubernetesSourceGitURL() string {
	s := os.Getenv("KUP_K8S_REPO_URL")
	if s != "" {
		return JoinSchema(s)
	}
	return JoinSchema(kubernetesSourceGitURL)
}

const kubernetesStableVersionURL = "dl.k8s.io/release/stable.txt"

func GetKubernetesStableVersionURL() string {
	s := os.Getenv("KUP_K8S_STABLE_URL")
	if s != "" {
		return JoinSchema(s)
	}
	return JoinSchema(kubernetesStableVersionURL)
}

const kubernetesDownloadBaseURL = "dl.k8s.io/release/%s/bin/%s/%s/kubectl" // version/os/arch

func GetKubernetesDownloadBaseURL(version string) string {
	s := os.Getenv("KUP_K8S_DL_URL")
	if s != "" {
		return JoinSchema(s)
	}
	return JoinSchema(fmt.Sprintf(kubernetesDownloadBaseURL, version, runtime.GOOS, runtime.GOARCH))
}

const DefaultKubernetesAcceleratedDownloadURL = "files.m.daocloud.io"

func GetKubernetesAcceleratedDownloadURL(version string) string {
	s := os.Getenv("KUP_K8S_MDL_URL")
	if s != "" {
		return JoinSchema(fmt.Sprintf("%s/%s", s, fmt.Sprintf(kubernetesDownloadBaseURL, version, runtime.GOOS, runtime.GOARCH)))
	}
	return JoinSchema(fmt.Sprintf("%s/%s", DefaultKubernetesAcceleratedDownloadURL, fmt.Sprintf(kubernetesDownloadBaseURL, version, runtime.GOOS, runtime.GOARCH)))
}

func GetBaseDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.Join(user.HomeDir, ".kube", "bin")
}

func GetCurrentKubectlPath() string {
	return filepath.Join(GetBaseDir(), "kubectl")
}

func GetKubectlPathWithVersion(version string) string {
	return filepath.Join(GetBaseDir(), fmt.Sprintf("kubectl-%s", version))
}

func CurrentKubectlIsSymLink() bool {
	stat, err := os.Lstat(GetCurrentKubectlPath())
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalln(err)
	}
	return stat.Mode()&os.ModeSymlink != 0
}

func KubectlIsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		log.Fatalln(err)
	}

	return true
}

func RemoveKubectl(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalln(err)
	}
}

func SetKubectlSysLink(version string) error {
	source := filepath.Join(GetKubectlPathWithVersion(version))
	target := filepath.Join(GetCurrentKubectlPath())

	if !KubectlIsExist(GetCurrentKubectlPath()) && !CurrentKubectlIsSymLink() {
		return os.Symlink(source, target)
	}

	if CurrentKubectlIsSymLink() {
		RemoveKubectl(GetCurrentKubectlPath())
		return os.Symlink(source, target)
	}

	return fmt.Errorf("Current kubectl file is not a symbolic link")
}
