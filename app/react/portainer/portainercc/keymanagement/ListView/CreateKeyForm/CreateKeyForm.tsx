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
import { useState } from 'react';

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

    const [file, setFile] = useState<File>();

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
                                            placeholder="Grant one or more teams access"
                                        />
                                    </FormControl>

                                    <FormControl
                                        inputId="key_import"
                                        label="Import"
                                    >

                                        <Field
                                            onChange={(event: any) => setFile(event.target.files[0])}
                                            type="File"
                                            name="import"
                                            id="key_import"
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
        console.log(values)
        //TODO
        let access = values.teamAccessPolicies.reduce((prev: any, current: any) => {
            return {
                ...prev,
                [current.toString()]: {
                    "RoleId": 0
                }
            }
        }, {})

        if (file) {
            let content: string = await readFileContent(file)
            const data = await createKey(keytype, values.description, access, content)
            console.log(data);
        } else {
            const data = await createKey(keytype, values.description, access)
            console.log(data);
        }
        
        return null;
    }

    function readFileContent(file: File): Promise<string> {
        return new Promise((resolve, reject) => {
            var fr = new FileReader();
            fr.onload = () => {
                resolve(fr.result as string);
            }
            fr.onerror = reject;
            fr.readAsText(file);
        })
    }

}
