import axios, { parseAxiosError } from '@/portainer/services/axios';
import { ServiceAdd, ServiceDeploy } from './types';



export async function addService(payload: ServiceAdd) {
    try {
        const { data } = await axios.post(buildUrl("add"), payload);
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

export async function deployService(payload: ServiceDeploy) {
    try {
        const { data } = await axios.post(buildUrl("deploy"), payload);
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

function buildUrl(id?: string, action?: string) {
    let url = '/ra/services';

    if (id) {
        url += `/${id}`;
    }

    if (action) {
        url += `/${action}`;
    }

    return url;
}