import { Formik, Field, Form } from 'formik';

import { PageHeader } from '@@/PageHeader';
import { CoordinatorListEntry } from '@/react/portainer/portainercc/coordinator/types';
import { Icon } from '@/react/components/Icon';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';

import { KeySelector } from '@@/KeySelector';


export function CoordinatorDeploymentView() {


    const coordintaorQuery = null;

    const exampleCoordinatorResult: CoordinatorListEntry[] = [
        {
            id: 1,
            name: "moin",
            imageId: "AF39BBAD222",
            signingKeyId: 1,
            uniqueId: "ABC123",
            signerId: "DEF999"
        },
        {
            id: 2,
            name: "cool",
            imageId: "AF39BBAD222",
            signingKeyId: 1,
            uniqueId: "ABC123",
            signerId: "DEF999"
        }
    ]

    let title = "Coordinator deployment";


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
                                initialValues={{key: 0}}
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


                                            <KeySelector
                                                value={values.key}
                                                onChange={(key) => setFieldValue('key', key)}
                                                keys={[]}
                                                placeholder="Select a coordinator image to deploy"
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
<hr/>
                                    Wenn einer vorhanden ist:
                                <hr/>

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
                                initialValues={{key: 0}}
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


        </>
    );
}
