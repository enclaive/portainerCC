import { Formik, Form } from 'formik';

import { PageHeader } from '@@/PageHeader';
import { Icon } from '@/react/components/Icon';
import clsx from 'clsx';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { Checkbox } from '@@/form-components/Checkbox';

import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { CoordinatorImageSelector } from '@@/CoordinatorImageSelector';
import { useCoordinatorImages } from '@/react/portainer/portainercc/coordinator/queries';
import { useCoordinatorDeploymentForEnv } from '../queries';
import { FormValues } from './types';
import { deployCoordinator, verifiyCoordinator } from '../coordinator.service';
import { CertInfoModalButton } from './CertInfoModal';

export function CoordinatorDeploymentView() {

    const envId = Number(useEnvironmentId());

    let deploymentQuery = useCoordinatorDeploymentForEnv(envId)
    const coordintaorQuery = useCoordinatorImages();

    let title = "Environment Coordinator";

    if (!coordintaorQuery.data) {
        return null;
    }

    const initialValues = {
        coordinatorImageId: 0,
        verify: true
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


                                            <FormControl
                                                inputId="verify"
                                                label="Auto verify after deployment?"
                                            >

                                                <Checkbox
                                                    id="verify"
                                                    label="Verify coordinator quote after deployment"
                                                    checked={values.verify}
                                                    onChange={() =>
                                                        setFieldValue('verify', !values.verify)
                                                    }
                                                />

                                            </FormControl>


                                            <div className="form-group">
                                                <div className="col-sm-12">
                                                    <LoadingButton
                                                        disabled={!isValid}
                                                        data-cy="team-createTeamButton"
                                                        isLoading={isSubmitting}
                                                        loadingText="Deploying coordinator, this may take a while..."
                                                    >
                                                        <Icon icon="plus" feather size="md" />
                                                        Deploy
                                                    </LoadingButton>
                                                </div>
                                            </div>
                                        </Form>
                                    )}
                                </Formik>

                                                    {/* <CertInfoModalButton /> */}
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
                                                            <pre style={{ overflow: "scroll", whiteSpace: "pre-wrap" }}>
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
        let pem = `-----BEGIN CERTIFICATE-----
        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
        -----END CERTIFICATE-----`
        console.log(pem);
        return null;
    }

    async function handleDeployment(values: FormValues) {
        await deployCoordinator(envId, values.coordinatorImageId)
        if (values.verify) {
            await new Promise(s => setTimeout(s, 10000));
            await verifiyCoordinator(envId)
        }
        return null;
    }
}
