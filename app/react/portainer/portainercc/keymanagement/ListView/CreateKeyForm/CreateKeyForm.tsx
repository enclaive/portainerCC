import { Formik, Field, Form } from 'formik';

import { Icon } from '@/react/components/Icon';

import { useState } from 'react';
import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { TeamsSelector } from '@@/TeamsSelector';
import { Team } from '../../../../users/teams/types'

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

    const [selectedTeams, setSelectedTeams] = useState<readonly number[]>([]);

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
                            initialValues={[]}
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
                                        inputId="key_teams"
                                        label="Teams"
                                        errors="err"
                                        required
                                    >
                                        <TeamsSelector
                                            value={selectedTeams}
                                            onChange={(value) => setSelectedTeams(value)}
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
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Creating team..."
                                                color='secondary'
                                            >
                                                <Icon icon="download" feather size="md" />
                                                Export
                                            </LoadingButton>
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Creating team..."
                                                color='danger'
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
