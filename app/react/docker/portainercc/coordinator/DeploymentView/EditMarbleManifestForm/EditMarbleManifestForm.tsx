import { Formik, Field, Form } from 'formik';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';

import { FormValues } from '@/react/portainer/portainercc/coordinator/ListView/BuildCoordinatorForm/types'
import { MarbleManifest } from '../types';

interface Props {
    manifest: MarbleManifest
}

export function EditMarbleManifestForm({ manifest }: Props) {
    let title = "Edit your manifest file"

    const initialValues = {
        name: '',
        key: 0
    }


    return (
        <div className="row">
            <div className="col-lg-12 col-md-12 col-xs-12">
                <Widget>
                    <Widget.Title
                        icon="edit"
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
                                            placeholder="e.g. super coordinator"
                                            data-cy="team-teamNameInput"
                                        />
                                    </FormControl>


                                {/* Buttons */}
                                    <div className="form-group">
                                        <div className="col-sm-12">
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Building coordinator image, this may take a while..."
                                                // onClick={() => handleBuildClick(values)}
                                            >
                                                Update Manifest
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
        return null;
    }

    // function handleBuildClick(values: any) {
    //     console.log("MOIN");
    //     console.log(values)
    // }

}