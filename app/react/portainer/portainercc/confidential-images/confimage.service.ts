import axios, { parseAxiosError } from '@/portainer/services/axios';

export async function getConfImages() {
    try {
        // CHANGE WHEN LIVE
        // const { data } = await axios.get<CoordinatorListEntry[]>(buildUrl("list"));
        // return data;
        return [
            {
                id: 1,
                imageid: "dockerID",
                mrsigner: "AF39BBAD222",
                mrenclave: "ABC123",
                timestamp: new Date(1668275068)
            },
            {
                id: 2,
                imageid: "dockerID",
                mrsigner: "AF39BBAD222",
                mrenclave: "ABC123",
                timestamp: new Date(1668225068)
            },
            {
                id: 3,
                imageid: "dockerID",
                mrsigner: "AF39BBAD222",
                mrenclave: "ABC123",
                timestamp: new Date(1668274068)
            }
        ]
    } catch (error) {
        throw parseAxiosError(error as Error);
    }
}

function buildUrl(id?: string, action?: string) {
    let url = '/ra/coordinator';

    if (id) {
        url += `/${id}`;
    }

    if (action) {
        url += `/${action}`;
    }

    return url;
}