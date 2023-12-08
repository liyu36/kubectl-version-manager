package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/liyu36/kubectl-version-manager/downloader"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd())
}

var (
	InstallCmdProgressFlag  bool
	InstallCmdOfficialFlag  bool
	InstallCmdReinstallFlag bool
)

func installCmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install the specified version of kubectl, default to the latest stable version",
		Long: fmt.Sprintf(`
	Install the specified version of the kubectl command,
	default to using the %s address to download the latest stable version`, GetKubernetesAcceleratedDownloadURL("vx.x.x")),
		Args:    cobra.MaximumNArgs(1),
		Example: fmt.Sprintf("  %s %s v1.28.4\n  %s %s // it install latest stable version", rootCmd.Name(), "install", rootCmd.Name(), "install"),
		Run:     install,
	}

	installCmd.Flags().BoolVarP(&InstallCmdProgressFlag, "quiet", "q", false, "Hide download progress")
	installCmd.Flags().BoolVarP(&InstallCmdOfficialFlag, "official", "o", false, "Download the kubectl command using the official address.")
	installCmd.Flags().BoolVarP(&InstallCmdReinstallFlag, "reinstall", "r", false, "Reinstall the specified version of kubectl")
	return installCmd
}

func install(cmd *cobra.Command, args []string) {
	var version string
	if len(args) > 0 {
		version = args[0]
	} else {
		version = GetKubernetesStableVersion()
	}

	if KubectlIsExist(GetKubectlPathWithVersion(version)) && !InstallCmdReinstallFlag {
		log.Printf("kubectl-%s is exist", version)
		return
	}

	if err := download(version); err != nil {
		log.Fatalln(err)
	}

	if err := os.Chmod(GetKubectlPathWithVersion(version), 0755); err != nil {
		log.Fatalln(err)
	}

	if err := SetKubectlSysLink(version); err != nil {
		log.Fatalln(err)
	}
}

func download(version string) error {
	var url = GetKubernetesAcceleratedDownloadURL(version)

	if InstallCmdOfficialFlag {
		url = GetKubernetesDownloadBaseURL(version)
	}

	log.Printf("Download from the %s\n", url)

	d, err := downloader.NewDownloader(url,
		downloader.DownloaderWithDirectory(GetBaseDir()),
		downloader.DownloaderWithFilename(fmt.Sprintf("kubectl-%s", version)),
		downloader.DownloaderWithProgress(!InstallCmdProgressFlag))
	if err != nil {
		return err
	}

	if err := d.Head(); err != nil {
		log.Println(err)
		return err
	}

	if err := d.Download().Merge(); err != nil {
		return err
	}
	fmt.Printf("kubectl-%s download complete\n", version)
	return nil
}

func GetKubernetesStableVersion() string {
	resp, err := http.Get(GetKubernetesStableVersionURL())
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
