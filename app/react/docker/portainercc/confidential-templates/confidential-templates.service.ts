import axios, { parseAxiosError } from '@/portainer/services/axios';
import { DeployConfidentialTemplateFormValues } from './DeploymentView/types';
import { ConfidentialTemplate } from './types';

export async function getConfidentialTemplates() {
    try {
        const { data } = await axios.get<ConfidentialTemplate[]>(buildUrl());
        return data;
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

export async function deployTemplate(payload: DeployConfidentialTemplateFormValues) {
    try {
        const { data } = await axios.post(buildUrl(), payload);
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

function buildUrl(id?: string, action?: string) {
    let url = '/portainercc/confidential-templates';

    if (id) {
        url += `/${id}`;
    }

    if (action) {
        url += `/${action}`;
    }

    return url;
}