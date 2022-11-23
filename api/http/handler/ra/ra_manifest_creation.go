package ra

import (
	"encoding/base64"

	portainer "github.com/portainer/portainer/api"
)

func createManifestMariadb(username, password, name string, initial bool) (portainer.CoordinatorManifest, map[string]map[string]string) {
	manifest := portainer.CoordinatorManifest{
		Packages: map[string]portainer.PackageProperties{
			name: {
				UniqueID: "e556b1f4a686be466b24c8d13df07705d5c1b9c8441281b84aec16c3d778521b",
			},
		},
		Marbles: map[string]portainer.Marble{
			name + "_marble": {
				Package: name,
				Parameters: portainer.Parameters{
					Files: map[string]portainer.File{
						"/app/init.sql": {
							Data:        "{{ raw .Secrets.init.Private }}",
							Encoding:    "string",
							NoTemplates: false,
						},
						"/dev/attestation/keys/default": {
							Data:        "{{ raw .Secrets.app_defaultkey.Private }}",
							Encoding:    "string",
							NoTemplates: false,
						},
					},
					Argv: []string{
						"/app/mariadbd",
						"--init-file=/app/init.sql",
					},
				},
			},
		},
		Secrets: map[string]portainer.Secret{
			"init": {
				Type:        "plain",
				UserDefined: true,
			},
		},
	}

	if initial {
		manifest.Secrets["app_defaultkey"] = portainer.Secret{
			Type: "symmetric-key",
			Size: 128,
		}
	}

	secretData := "CREATE OR REPLACE USER " + username + " IDENTIFIED BY '" + password + "';\n GRANT ALL PRIVILEGES ON *.* TO " + username + ";"
	secretBase64 := base64.StdEncoding.EncodeToString([]byte(secretData))

	secretMap := map[string]map[string]string{}
	secretMap["init"] = make(map[string]string)
	secretMap["init"]["Key"] = secretBase64

	return manifest, secretMap
}
