import axios, { parseAxiosError } from '@/portainer/services/axios';
import { CoordinatorDeployment } from './types';

export async function getCoordinatorDeployments() {
    try {
        const { data } = await axios.get<CoordinatorDeployment[]>(buildUrl("list"));
        return data;
        // return [{
        //     id: 22,
        //     coordinatorId: 4,
        //     endpointId: 11,
        //     rootCert: "string",
        //     userCert: "string",
        //     userPrivKey: "string",
        //     manifest: "any",
        //     verified: false
        // }]
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

export async function getCoordinatorDeploymentForEnv(envId: number) {
    let all = await getCoordinatorDeployments();
    return all.find(dep => dep.endpointId == envId);
}

export async function deployCoordinator(envId: number, coordinatorId: number){
    try {
        const { data } = await axios.post(buildUrl(), { environmentId: envId, coordinatorId: coordinatorId });
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

export async function verifiyCoordinator(envId: number) {
    try {
        const { data } = await axios.get<CoordinatorDeployment[]>("/ra/coordinator/verify/" + envId.toString());
        return data;
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

function buildUrl(id?: string, action?: string) {
    let url = '/ra/coordinator/deploy';

    if (id) {
        url += `/${id}`;
    }

    if (action) {
        url += `/${action}`;
    }

    return url;
}