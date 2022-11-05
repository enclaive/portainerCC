import { Formik, Field, Form } from 'formik';

import { Icon } from '@/react/components/Icon';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { KeyEntry } from '@/react/portainer/portainercc/keymanagement/types';

import { KeySelector } from '@@/KeySelector';
import { FormValues } from '@/react/portainer/portainercc/coordinator/ListView/BuildCoordinatorForm/types'
import { buildCoordinator } from '../../coordinator.service';

interface Props {
    keys: KeyEntry[]
}

export function BuildCoordinatorForm({ keys }: Props) {
    let title = "Build a new coordinator image"

    const initialValues = {
        name: '',
        key: 0
    }


    return (
        <div className="row">
            <div className="col-lg-12 col-md-12 col-xs-12">
                <Widget>
                    <Widget.Title
                        icon="plus"
                        title={title}
                        featherIcon
                        className="vertical-center"
                    />
                    <Widget.Body>
                        <Formik
                            initialValues={initialValues}
                            onSubmit={handleBuild}
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
                                        inputId="key_name"
                                        label="Name"
                                        errors="err"
                                        required
                                    >
                                        <Field
                                            as={Input}
                                            name="name"
                                            id="key_name"
                                            required
                                            placeholder="e.g. super key"
                                            data-cy="team-teamNameInput"
                                        />
                                    </FormControl>


                                    <FormControl
                                        inputId="key"
                                        label="SGX Signign Key"
                                        errors="err"
                                        required
                                    >


                                        <KeySelector
                                            value={values.key}
                                            onChange={(key) => setFieldValue('key', key)}
                                            keys={keys}
                                            placeholder="Select a key"
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
                                                Build
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
    );

    async function handleBuild(values: FormValues) {
        const data = await buildCoordinator(values.name, values.key)
        console.log(data);
        return null;
    }

    // function handleBuildClick(values: any) {
    //     console.log("MOIN");
    //     console.log(values)
    // }

}