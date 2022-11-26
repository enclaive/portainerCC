import { Formik, Form } from 'formik';

import { PageHeader } from '@@/PageHeader';
import { Icon } from '@/react/components/Icon';
import clsx from 'clsx';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';

import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { CoordinatorImageSelector } from '@@/CoordinatorImageSelector';
import { useCoordinatorImages } from '@/react/portainer/portainercc/coordinator/queries';
import { useCoordinatorDeploymentForEnv } from '../queries';
import { FormValues } from './types';
import { deployCoordinator, verifiyCoordinator } from '../coordinator.service';

export function CoordinatorDeploymentView() {

    const envId = Number(useEnvironmentId());

    let deploymentQuery = useCoordinatorDeploymentForEnv(envId)
    const coordintaorQuery = useCoordinatorImages();

    let title = "Environment Coordinator";

    if (!coordintaorQuery.data) {
        return null;
    }

    const initialValues = {
        coordinatorImageId: 0
    }

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {!deploymentQuery.data &&
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
                                    initialValues={initialValues}
                                    onSubmit={handleDeployment}
                                    key={1}
                                >
                                    {({
                                        values,
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
                                                required
                                            >

                                                <CoordinatorImageSelector
                                                    value={values.coordinatorImageId}
                                                    onChange={(coordinatorImageId) => setFieldValue('coordinatorImageId', coordinatorImageId)}
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
                                                        loadingText="Deploying coordinator, this may take a while..."
                                                    // onClick={() => handleBuildClick(values)}
                                                    >
                                                        <Icon icon="plus" feather size="md" />
                                                        Deploy and Verify
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
            }


            {deploymentQuery.data &&
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
                                    onSubmit={handleVerifyClick}
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

                                            <table className='table'>
                                                <tbody>
                                                    <tr>
                                                        <td className="col-xs-6 col-sm-4 col-md-3 col-lg-3">Coordinator ID</td>
                                                        <td>{deploymentQuery.data?.coordinatorId}</td>
                                                    </tr>
                                                    <tr>
                                                        <td>Status</td>
                                                        <td>
                                                            {deploymentQuery.data?.verified && <>
                                                                <Icon
                                                                    icon="lock"
                                                                    className={clsx('icon icon-sm icon-success')}
                                                                    feather
                                                                />
                                                                Verified
                                                            </>
                                                            }
                                                            {!deploymentQuery.data?.verified &&
                                                                <>
                                                                    <Icon
                                                                        icon="lock"
                                                                        className={clsx('icon icon-sm icon-danger')}
                                                                        feather
                                                                    />
                                                                    Not Verified
                                                                </>
                                                            }
                                                        </td>
                                                    </tr>
                                                    <tr>
                                                        <td>Manifest:</td>
                                                        <td>
                                                            <pre style={{overflow: "scroll", whiteSpace: "pre-wrap"}}>
                                                                {JSON.stringify(deploymentQuery.data?.manifest, null, 2)}
                                                            </pre>
                                                        </td>
                                                    </tr>

                                                </tbody>
                                            </table>

                                            <div className="form-group">
                                                <div className="col-sm-12">
                                                    <LoadingButton
                                                        disabled={!isValid}
                                                        isLoading={isSubmitting}
                                                        loadingText="Verifiying coordinator deployment..."
                                                    >
                                                        <Icon icon="shield" feather size="md" />
                                                        Verify
                                                    </LoadingButton>
                                                    {/* <LoadingButton
                                                        disabled={!deploymentQuery.data?.rootCert.Bytes}
                                                        isLoading={isSubmitting}
                                                        loadingText="Building coordinator image, this may take a while..."
                                                        type='button'
                                                        onClick={handleCertClick}
                                                    >
                                                        <Icon icon="lock" feather size="md" />
                                                        Inspect Root Certificate
                                                    </LoadingButton>
                                                    <LoadingButton
                                                        disabled={!deploymentQuery.data?.rootCert.Bytes}
                                                        isLoading={isSubmitting}
                                                        loadingText="Building coordinator image, this may take a while..."
                                                        type='button'
                                                        onClick={handleCertClick}
                                                    >
                                                        <Icon icon="lock" feather size="md" />
                                                        Inspect User Certificate
                                                    </LoadingButton> */}
                                                </div>
                                            </div>
                                        </Form>
                                    )}
                                </Formik>


                            </Widget.Body>
                        </Widget>
                    </div>
                </div >
            }
        </>
    );

    async function handleVerifyClick() {
        const data = await verifiyCoordinator(envId)
        console.log(data);
        return null;
    }

    async function handleCertClick() {
        console.log("cert")
        return null;
    }

    async function handleDeployment(values: FormValues) {
        const data = await deployCoordinator(envId, values.coordinatorImageId)
        console.log(data)
        return null;
    }
}
