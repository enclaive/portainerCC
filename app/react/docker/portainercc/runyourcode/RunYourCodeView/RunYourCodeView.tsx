import { PageHeader } from '@@/PageHeader';
import { Formik, Field, Form, FieldArray } from 'formik';

import { Input, Select } from '@@/form-components/Input';
import { FormControl } from '@@/form-components/FormControl';
import { Widget } from '@@/Widget';
import { LoadingButton } from '@@/buttons/LoadingButton';
import { FormSection } from '@@/form-components/FormSection';
import { KeySelector } from '@@/KeySelector';
import { useKeys } from '@/react/portainer/portainercc/keymanagement/queries';

import { useEnvironmentId } from '@/portainer/hooks/useEnvironmentId';
import { FormValues } from './types';
import { Button } from '@@/buttons';
import { run } from '../runyourcode.service';

export function RunYourCodeView() {
  let title = 'RunYourCode';

  let env = useEnvironmentId();
  const envId = Number(env);


  const keysQuery = useKeys('SIGNING')
  if (!keysQuery.data || !env) {
    return null;
  }

  const initialValues: FormValues = {
    Type: "node",
    EnvId: envId,
    SigningKeyId: 0,
    Name: "",
    Ports: [],
    Repository: "",
    BuildArgs: "",
    RunArgs: ""
  }

  return (
    <>
      <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

      <div className="row">
        <div className="col-lg-12 col-md-12 col-xs-12">
          <Widget>
            <Widget.Title
              icon="command"
              title={title}
              featherIcon
              className="vertical-center"
            />
            <Widget.Body>
              <Formik
                enableReinitialize
                initialValues={initialValues}
                onSubmit={handleSubmit}
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

                    <FormSection title='Image Setup'>

                      <FormControl inputId="Type" label="Language" required>
                        <Field as={Select} name="Type" id="Type" options={[{ value: "node", label: "Node.js" }, { value: "python", label: "Python3" }]} />
                      </FormControl>

                      <FormControl inputId="key" label="SGX Signign Key" required>
                        <KeySelector
                          value={values.SigningKeyId}
                          onChange={(key) => setFieldValue('SigningKeyId', key)}
                          keys={keysQuery.data}
                          placeholder="Select a key"
                        />
                      </FormControl>

                      <FormControl inputId="Name" label="Name" required>
                        <Field as={Input} name="Name" id="Name" required />
                      </FormControl>

                      <FormControl inputId="Repository" label="Repository URL" required>
                        <Field as={Input} name="Repository" id="Repository" required />
                      </FormControl>

                      <FormControl inputId="RunArgs" label="RunArg" required>
                        <Field as={Input} name="RunArgs" id="RunArgs" required />
                      </FormControl>

                      <div>
                        <span className='small text-muted'>Your Files (cloned Github Repository) will be available under /app/. (e.g. RunArg: /app/server.js)</span>
                      </div>

                      <FormSection title='Port mapping'>
                        <FieldArray name="Ports" render={(arrayHelpers) => (
                          <div className='p-10'>
                            {values.Ports.map((entry, index) => (
                              <div className="form-group" key={index}>
                                <div className="form-inline">
                                  <div className="input-group col-sm-3">
                                    <span className="input-group-addon">Type</span>
                                    <Field as={Select} name={`Ports.${index}.Type`} id={`Ports.${index}.Type`} options={[{ value: "tcp", label: "TCP" }, { value: "udp", label: "UDP" }]} />
                                  </div>

                                  <div className="input-group col-sm-4 input-group-sm">
                                    <span className="input-group-addon">Host</span>
                                    <Field as={Input} name={`Ports.${index}.Host`} id={`Ports.${index}.Host`} required />
                                  </div>

                                  <div className="input-group col-sm-4 input-group-sm">
                                    <span className="input-group-addon">Container</span>
                                    <Field as={Input} name={`Ports.${index}.Container`} id={`Ports.${index}.Container`} required />
                                  </div>

                                  <div className="input-group col-sm-1">
                                    <Field as={Button} name="RmPort" id="RmPort" onClick={() => arrayHelpers.remove(index)}>Remove</Field>
                                  </div>

                                </div>
                              </div>
                            ))}
                            <Field as={Button} name="AddPort" id="AddPort" onClick={() => arrayHelpers.push({ Type: "tcp", Host: "", Container: "" })}>Add Port</Field>
                          </div>
                        )} />
                      </FormSection>

                    </FormSection>


                    <div className="form-group">
                      <div className="col-sm-12">
                        <LoadingButton
                          disabled={!isValid}
                          isLoading={isSubmitting}
                          loadingText="Adding a confidential service..."
                        >
                          Build and Deploy your Github Repository
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

  async function handleSubmit(values: FormValues) {
    console.log(values)
    const res = await run(values);

    console.log(res);
    return;
  }

}
