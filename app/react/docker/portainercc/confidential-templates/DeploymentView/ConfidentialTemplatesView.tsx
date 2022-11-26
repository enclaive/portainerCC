import { Formik, Field, Form } from 'formik';

import { PageHeader } from '@@/PageHeader';
import { Input } from '@@/form-components/Input';


import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { FormValues } from '@/react/docker/portainercc/confidential-templates/DeploymentView/types'
import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { addService, deployService } from '../service-deployment.service';

export function ConfidentialTemplatesView() {


    let title = "Confidential Templates";

    const envId = Number(useEnvironmentId());

    const initialValues = {
        Name: '',
        Username: '',
        Password: '',
        ImageID: 'sgxdcaprastuff/gramine-mariadb'
    }


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            <div className="row">
                <div className="col-lg-12 col-md-12 col-xs-12">
                    <Widget>
                        <Widget.Title
                            icon="codesandbox"
                            title="MariaDB Example"
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
                                            inputId="ImageID"
                                            label="ImageID"
                                            required
                                        >
                                            <Field
                                                as={Input}
                                                name="ImageID"
                                                id="ImageID"
                                                required

                                            />
                                        </FormControl>

                                        <FormControl
                                            inputId="Name"
                                            label="Name"
                                            required
                                        >
                                            <Field
                                                as={Input}
                                                name="Name"
                                                id="Name"
                                                required
                                                placeholder="e.g. mariadb"
                                            />
                                        </FormControl>
                                        <FormControl
                                            inputId="Username"
                                            label="Username"
                                            required
                                        >
                                            <Field
                                                as={Input}
                                                name="Username"
                                                id="Username"
                                                required
                                                placeholder="e.g. root"
                                            />
                                        </FormControl>
                                        <FormControl
                                            inputId="Password"
                                            label="Password"
                                            required
                                        >
                                            <Field
                                                as={Input}
                                                name="Password"
                                                id="Password"
                                                required
                                                placeholder="e.g. secret"
                                                type="password"
                                            />
                                        </FormControl>


                                        <div className="form-group">
                                            <div className="col-sm-12">
                                                <LoadingButton
                                                    disabled={!isValid}
                                                    data-cy="team-createTeamButton"
                                                    isLoading={isSubmitting}
                                                    loadingText="Adding a confidential service..."
                                                >
                                                    Add & Deploy Service
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


    async function handleDeployment(values: FormValues) {
        console.log("MOIN")
        const addResult = await addService({ EnvironmentID: envId, Name: values.Name, Username: values.Username, Password: values.Password })
        console.log(addResult)
        const deployResult = await deployService({ EnvironmentID: envId, Name: values.Name, ImageID: values.ImageID })
        console.log(deployResult);
        return;
    }

}
