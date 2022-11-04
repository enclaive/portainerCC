import { Formik, Field, Form } from 'formik';

import { Icon } from '@/react/components/Icon';

import { useState } from 'react';
import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { KeyEntry } from '../../../keymanagement/types';
import { KeySelector } from '@@/KeySelector';

interface Props {
    keys: KeyEntry[]
}

export function BuildCoordinatorForm({
    keys
}: Props) {
    let title = "Build your coordinator image"

    const [selectedKey, setSelectedKey] = useState<readonly number[]>([]);


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
                            initialValues={{ name: '', key: null }}
                            onSubmit={(() => Promise.resolve(null))}
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
                                        inputId="keys"
                                        label="SGX Signign Key"
                                        errors="err"
                                        required
                                    >


                                        <KeySelector
                                            value={selectedKey}
                                            onChange={(keys) => setSelectedKey(keys)}
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
                                                loadingText="Creating key..."
                                            >
                                                <Icon icon="plus" feather size="md" />
                                                Build
                                            </LoadingButton>
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Creating team..."
                                            >
                                                <Icon icon="trash" feather size="md" />
                                                Remove
                                            </LoadingButton>
                                        </div>
                                    </div>
                                </Form>
                            )}
                        </Formik>
                    </Widget.Body>
                </Widget>
            </div>
        </div>
    );
}