import axios, { parseAxiosError } from '@/portainer/services/axios';
import { ConfidentialTemplate } from './types';
import { FormValues } from './RunYourCodeView/types';

export async function getConfidentialTemplates() {
    try {
        const { data } = await axios.get<ConfidentialTemplate[]>(buildUrl());
        return data;
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

export async function run(values: FormValues) {
    try {

        const payload = {
            EnvId: values.EnvId,
            SigningKeyId: values.SigningKeyId,
            Name: values.Name,
            Ports: values.Ports,
            Repository: values.Repository,
            BuildArgs: "",
            RunArgs: values.RunArgs
        }

        const { data } = await axios.post(buildUrl(values.Type), payload);
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

function buildUrl(type?: string) {
    let url = '/portainercc/confidential-templates';

    if (type) {
        url += `/${type}`;
    }

    return url;
}