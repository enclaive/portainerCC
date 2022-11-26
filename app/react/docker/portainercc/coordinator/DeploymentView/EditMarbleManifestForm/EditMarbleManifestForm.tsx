import { useState } from 'react';
import { Formik, Field, Form } from 'formik';

import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { Input } from '@@/form-components/Input';
import { LoadingButton } from '@@/buttons/LoadingButton';

import { MarbleManifest } from '../types';

interface Props {
    origManifest: MarbleManifest
}

export function EditMarbleManifestForm({ origManifest }: Props) {
    let title = "Edit your manifest file"

    const [manifest, setManifest] = useState(origManifest);

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
                            initialValues={manifest}
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
                                    <label>Hallo</label>

                                    {/* Buttons */}
                                    <div className="form-group">
                                        <div className="col-sm-12">
                                            <LoadingButton
                                                disabled={!isValid}
                                                data-cy="team-createTeamButton"
                                                isLoading={isSubmitting}
                                                loadingText="Building coordinator image, this may take a while..."
                                                type='button'
                                                onClick={() => test()}
                                            >
                                                TEST
                                            </LoadingButton>
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

    async function test() {
        //add package
        manifest.Packages["addedPackage"] = {
            UniqueID: "123"
        }

        manifest.Marbles["addedMarble"] = {
            test: "test"
        }
        console.log("TEST")
        return null;
    }

    async function handleBuild(values: any) {
        console.log(manifest)
        return null;
    }

    // function handleBuildClick(values: any) {
    //     console.log("MOIN");
    //     console.log(values)
    // }

}