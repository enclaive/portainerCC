import { Formik, Field, Form } from 'formik';

import { Icon } from '@/react/components/Icon';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { TeamsSelector } from '@@/TeamsSelector';
import { Team } from '../../../../users/teams/types'
import { FormValues } from './types';
import { createKey } from '../../keys.service';

interface Props {
    keytype: string;
    teams: Team[];
}

export function CreateKeyForm({ keytype, teams }: Props) {
    let title = ""
    if (keytype == "SIGNING") {
        title = "Create or import a SGX signing key"
    } else if (keytype == "FILE_ENC") {
        title = "Create or import a gramine proteced files key"
    } else {
        throw Error("invalid key type")
    }

    const initialValues = {
        description: '',
        teamAccessPolicies: []
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
                            onSubmit={handleSubmit}
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
                                        inputId="key_description"
                                        label="Description"
                                        required
                                    >
                                        <Field
                                            as={Input}
                                            name="description"
                                            id="key_description"
                                            required
                                            placeholder="e.g. super key"
                                            data-cy="team-teamNameInput"
                                        />
                                    </FormControl>

                                    <FormControl
                                        inputId="key_teams"
                                        label="Teams"
                                        required
                                    >

                                        <TeamsSelector
                                            value={values.teamAccessPolicies}
                                            onChange={(values) => setFieldValue('teamAccessPolicies', values)}
                                            teams={teams}
                                            placeholder="Select one or more teams to access the key"
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
                                                Create
                                            </LoadingButton>
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Creating team..."
                                            >
                                                <Icon icon="upload" feather size="md" />
                                                Import
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

    async function handleSubmit(values: FormValues) {
        //TODO
        let access = values.teamAccessPolicies.reduce((prev: any, current: any) => {
            return {
                ...prev,
                [current.toString()]: {
                    "RoleId": 0
                }
            }
        }, {})
        const data = await createKey(keytype, values.description, access)
        console.log(data);
        return null;
    }
}
