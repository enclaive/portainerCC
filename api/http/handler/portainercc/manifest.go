package portainercc

import portainer "github.com/portainer/portainer/api"

func createUpdateManifest(template portainer.ConfidentialTemplate, inputParam ConfTempDeployParams, mrenclave string, mrsigner string) portainer.CoordinatorManifest {
	manifest := portainer.CoordinatorManifest{
		Packages: map[string]portainer.PackageProperties{
			inputParam.Name: {
				//mrenclave for now
				UniqueID: mrenclave,
			},
		},
		Marbles: map[string]portainer.Marble{
			inputParam.Name + "_marble": {
				Package:    inputParam.Name,
				Parameters: template.ManifestBoilerplate.ManifestParameters,
			},
		},
		Secrets: template.ManifestBoilerplate.ManifestSecrets,
	}

	return manifest
}
