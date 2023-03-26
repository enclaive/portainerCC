import axios, { parseAxiosError } from '@/portainer/services/axios';

export async function getEncryptedVolumes(endPointId: string) {
    try {
        const { data } = await axios.get<any[]>(buildUrl(endPointId));
        return data;
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

function buildUrl(endPointId?: string) {
    let url = '/endpoints/' + endPointId + '/docker/volumes?filters={"label":["encrypted=true"]}';

    return url;
}