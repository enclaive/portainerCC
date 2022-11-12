import { Formik, Field, Form } from 'formik';

import { PageHeader } from '@@/PageHeader';
import { CoordinatorListEntry } from '@/react/portainer/portainercc/coordinator/types';
import { Icon } from '@/react/components/Icon';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';

import { CoordinatorImageSelector } from '@@/CoordinatorImageSelector';
import { useCoordinatorImages } from '@/react/portainer/portainercc/coordinator/queries';
import { EditMarbleManifestForm } from './EditMarbleManifestForm/EditMarbleManifestForm';
import { MarbleManifest } from './types';

const exampleManifest: MarbleManifest = {
    "Packages": {
        "backend": {
            "UniqueID": "6b2822ac2585040d4b9397675d54977a71ef292ab5b3c0a6acceca26074ae585",
            "Debug": false
        },
        "frontend": {
            "SignerID": "43361affedeb75affee9baec7e054a5e14883213e5a121b67d74a0e12e9d2b7a",
            "ProductID": 43,
            "SecurityVersion": 3,
            "Debug": true
        }
    },
    "Marbles": {
        "backendFirst": {
            "Package": "backend",
            "MaxActivations": 1,
            "Parameters": {
                "Files": {
                    "/tmp/defg.txt": "foo",
                    "/tmp/jkl.mno": "bar",
                    "/tmp/pqr.ust": {
                        "Data": "Zm9vCmJhcg==",
                        "Encoding": "base64",
                        "NoTemplates": true
                    }
                },
                "Env": {
                    "IS_FIRST": "true",
                    "ROOT_CA": "{{ pem .MarbleRun.RootCA.Cert }}",
                    "MARBLE_CERT": "{{ pem .MarbleRun.MarbleCert.Cert }}",
                    "MARBLE_KEY": "{{ pem .MarbleRun.MarbleCert.Private }}"
                },
                "Argv": [
                    "--first",
                    "serve"
                ]
            },
            "TLS": [
                "backendFirstTLS"
            ]
        },
        "frontend": {
            "Package": "frontend",
            "Parameters": {
                "Env": {
                    "ROOT_CA": "{{ pem .MarbleRun.RootCA.Cert }}",
                    "MARBLE_CERT": "{{ pem .MarbleRun.MarbleCert.Cert }}",
                    "MARBLE_KEY": "{{ pem .MarbleRun.MarbleCert.Private }}"
                }
            },
            "TLS": [
                "frontendTLS1", "frontendTLS2"
            ]
        }
    }
}

export function CoordinatorDeploymentView() {


    const coordintaorQuery = useCoordinatorImages();

    let title = "Coordinator deployment";

    if (!coordintaorQuery.data) {
        return null;
    }

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            <div className="row">
                <div className="col-lg-12 col-md-12 col-xs-12">
                    <Widget>
                        <Widget.Title
                            icon="codesandbox"
                            title={title}
                            featherIcon
                            className="vertical-center"
                        />
                        <Widget.Body>
                            <Formik
                                initialValues={{ coordinator: 0 }}
                                onSubmit={() => Promise.resolve(null)}
                                key={1}
                            >
                                {({
                                    values,
                                    errors,
                                    handleSubmit,
                                    setFieldValue,
                                    isSubmitting,
                                    isValid,
                                }) => (
                                    <Form
                                        className="form-horizontal"
                                        onSubmit={handleSubmit}
                                        noValidate
                                    >

                                        <FormControl
                                            inputId="coordinatorId"
                                            label="Coordinator Image"
                                            errors="err"
                                            required
                                        >

                                            <CoordinatorImageSelector
                                                value={values.coordinator}
                                                onChange={(coordinator) => setFieldValue('coordinator', coordinator)}
                                                images={coordintaorQuery.data}
                                                placeholder="Select a coordinator image to deply"
                                            />

                                        </FormControl>


                                        <div className="form-group">
                                            <div className="col-sm-12">
                                                <LoadingButton
                                                    disabled={!isValid}
                                                    data-cy="team-createTeamButton"
                                                    isLoading={isSubmitting}
                                                    loadingText="Building coordinator image, this may take a while..."
                                                // onClick={() => handleBuildClick(values)}
                                                >
                                                    <Icon icon="plus" feather size="md" />
                                                    Deploy
                                                </LoadingButton>
                                            </div>
                                        </div>
                                    </Form>
                                )}
                            </Formik>


                        </Widget.Body>
                    </Widget>
                </div>
            </div >
            <hr />
            Wenn einer vorhanden ist:
            <hr />

            <div className="row">
                <div className="col-lg-12 col-md-12 col-xs-12">
                    <Widget>
                        <Widget.Title
                            icon="codesandbox"
                            title={title}
                            featherIcon
                            className="vertical-center"
                        />
                        <Widget.Body>
                            <Formik
                                initialValues={{ key: 0 }}
                                onSubmit={() => Promise.resolve(null)}
                                key={1}
                            >
                                {({
                                    values,
                                    errors,
                                    handleSubmit,
                                    setFieldValue,
                                    isSubmitting,
                                    isValid,
                                }) => (
                                    <Form
                                        className="form-horizontal"
                                        onSubmit={handleSubmit}
                                        noValidate
                                    >

                                        <div className="form-group">
                                            <div className='py-5'>Coordinator COORDINATOR_NAME is running in this environment.</div>
                                            <div className="col-sm-12">
                                                <LoadingButton
                                                    disabled={!isValid}
                                                    data-cy="team-createTeamButton"
                                                    isLoading={isSubmitting}
                                                    loadingText="Building coordinator image, this may take a while..."
                                                // onClick={() => handleBuildClick(values)}
                                                >
                                                    <Icon icon="shield" feather size="md" />
                                                    Reattest
                                                </LoadingButton>
                                                <LoadingButton
                                                    disabled={!isValid}
                                                    data-cy="team-createTeamButton"
                                                    isLoading={isSubmitting}
                                                    loadingText="Building coordinator image, this may take a while..."
                                                // onClick={() => handleBuildClick(values)}
                                                >
                                                    <Icon icon="lock" feather size="md" />
                                                    Inspect Coordinator Certificate
                                                </LoadingButton>
                                            </div>
                                        </div>
                                    </Form>
                                )}
                            </Formik>


                        </Widget.Body>
                    </Widget>
                </div>
            </div >

            <EditMarbleManifestForm manifest={exampleManifest} />

        </>
    );
}
