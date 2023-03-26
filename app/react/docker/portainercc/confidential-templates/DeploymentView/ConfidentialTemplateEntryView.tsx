import { Formik, Field, Form } from 'formik';

import { Input, Select } from '@@/form-components/Input';
import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { DeployConfidentialTemplateFormValues } from '@/react/docker/portainercc/confidential-templates/DeploymentView/types'
import { ConfidentialTemplate } from '../types';
import { FormSection } from '@@/form-components/FormSection';
import { deployTemplate } from '../confidential-templates.service';
import { useState } from 'react';

interface Props {
    template: ConfidentialTemplate
    envId: number
    encryptedVolumes: any //TODO any
}

export function ConfidentialTemplateEntryView({ template, envId, encryptedVolumes }: Props) {

    const [toggle, setToggle] = useState(false);

    const initialValues = {
        Id: template.Id,
        EnvId: envId,
        Image: template.Image,
        Name: "",
        Inputs: template.Inputs.reduce((acc, curr) => ({ ...acc, [curr.Label]: "" }), {})
    }

    let encVolOptions =
        1 > 0 &&
        encryptedVolumes.map((vol: any) => {
            return { value: vol.Name, label: vol.Name }
        });

    let secrets = template.Inputs.filter(input => input.Type == "SECRET");
    let ports = template.Inputs.filter(input => input.Type == "PORT");
    let volumes = template.Inputs.filter(input => input.Type == "VOLUME")

    return (
        <>
            <div className="row">
                <div className="col-lg-12 col-md-12 col-xs-12">
                    <Widget>
                        <div className="widget-header" onClick={() => setToggle(!toggle)}>
                            <div className="row">
                                <span className={'pull-left vertical-center'}>
                                    <div className='vertical-center justify-center min-w-[56px]'>
                                        <img className="blocklist-item-logo" src={template.LogoURL} />
                                    </div>
                                    <div className='blocklist-item-line'>
                                        <span className='ml-5 blocklist-item-title'>{template.TemplateName}</span>
                                    </div>
                                    {/* <div className='blocklist-item-line template-item-details-sub'>
                                        <span className='blocklist-item-desc'>Super sicherer container aus super sicherem image</span>
                                        <span className='small text-muted'>Bereitgestellt von marcel</span>
                                    </div> */}
                                </span>
                            </div>
                        </div>
                        {toggle && (<Widget.Body>
                            <Formik
                                initialValues={initialValues}
                                onSubmit={handleDeployment}
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

                                        <FormSection title='Package info'>
                                            <FormControl inputId="Image" label="Image">
                                                <div className='input-group'>
                                                    <Field readOnly as={Input} name="Image" id="Image" />
                                                    <span className='input-group-btn'>
                                                        <a
                                                            href={"https://hub.docker.com/r/" + template.Image}
                                                            className="btn btn-default vertical-center"
                                                            title="Show on Docker Hub"
                                                            target="_blank"
                                                        >
                                                            Show on Docker Hub
                                                        </a>
                                                    </span>
                                                </div>
                                            </FormControl>

                                            <FormControl inputId="Name" label="Name" required>
                                                <Field as={Input} name="Name" id="Name" required />
                                            </FormControl>
                                        </FormSection>

                                        {secrets.length > 0 && <FormSection title='Secrets'>
                                            {secrets.map((input) => {
                                                let str = "Inputs." + input.Label;
                                                return (
                                                    <>
                                                        <FormControl inputId={str} label={input.Label} required>
                                                            <Field as={Input} name={str} id={str} required placeholder={input.Default} />
                                                        </FormControl>
                                                    </>
                                                )
                                            })}
                                        </FormSection>}

                                        {ports.length > 0 && <FormSection title='Port mapping'>
                                            {ports.map((input) => {
                                                let str = "Inputs." + input.Label;
                                                return (
                                                    <>
                                                        <div className="form-group">
                                                            <div className="form-inline mx-10">
                                                                <div className="input-group col-sm-1">
                                                                    <div className="btn-group btn-group-sm">
                                                                        {input.Label}
                                                                    </div>
                                                                </div>
                                                                <div className="input-group col-sm-5 input-group-sm">
                                                                    <span className="input-group-addon">host</span>
                                                                    <Field as={Input} name={str} id={str} required placeholder={input.Default} />
                                                                </div>
                                                                <div className="input-group col-sm-5 input-group-sm">
                                                                    <span className="input-group-addon">container</span>
                                                                    <input type="text" className="form-control" readOnly value={input.PortContainer} />
                                                                </div>
                                                                <div className="input-group col-sm-1 input-group-sm">
                                                                    <div className="btn-group btn-group-sm">
                                                                        <label className="btn btn-light">{input.PortType}</label>
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </>
                                                )
                                            })}

                                        </FormSection>}

                                        {volumes.length > 0 && <FormSection title='Encrypted Volumes'>
                                            {volumes.map((input) => {
                                                let str = "Inputs." + input.Label;
                                                return (
                                                    <>
                                                        <div className="form-group">
                                                            <div className="col-sm-10 mx-10 input-group">
                                                                <span className="input-group-addon">{input.Label}</span>
                                                                <div className="col-sm-12 input-group">
                                                                    <Field as={Select} name={str} id={str} options={encVolOptions} />
                                                                    {/* <option key="create" value="create">Create new encrypted volume</option> */}
                                                                    {/* {encVolOptions} */}
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </>
                                                )
                                            })}

                                        </FormSection>}



                                        <div className="form-group">
                                            <div className="col-sm-12">
                                                <LoadingButton
                                                    disabled={!isValid}
                                                    isLoading={isSubmitting}
                                                    loadingText="Adding a confidential service..."
                                                >
                                                    Deploy Service
                                                </LoadingButton>
                                            </div>
                                        </div>
                                    </Form>
                                )}
                            </Formik>


                        </Widget.Body>)}
                    </Widget>
                </div>
            </div >
        </>
    );

    async function handleDeployment(values: DeployConfidentialTemplateFormValues) {
        console.log("MOIN")
        console.log(values)
        const deploy = await deployTemplate(values);

        console.log(deploy);
        // const addResult = await addService({ EnvironmentID: envId, Name: values.Name, Username: values.Username, Password: values.Password })
        // console.log(addResult)
        // const deployResult = await deployService({ EnvironmentID: envId, Name: values.Name, ImageID: values.ImageID })
        // console.log(deployResult);
        return;
    }

}
