import axios, { parseAxiosError } from '@/portainer/services/axios';
import { KeyId } from '../keymanagement/types';
import { CoordinatorImageId, CoordinatorListEntry } from './types';

// export async function getKeys(type: string) {
//     try {
//       const { data } = await axios.get<KeyEntry[]>(buildUrl(), {
//         params: { type },
//       });
//       return data;
//     } catch (error) {
//       throw parseAxiosError(error as Error);
//     }
//   }

//   function buildUrl(id?: KeyId) {
//     let url = '/portainercc/keys';

//     if (id) {
//       url += `/${id}`;
//     }

//     return url;
//   }

export async function buildCoordinator(name: string, keyId: KeyId) {
    try {
        const { data } = await axios.post(buildUrl(undefined, "build"), { Name: name, SigningKeyId: keyId });
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

export async function removeCoordinatorImage(id: CoordinatorImageId) {
    try {
        const { data } = await axios.delete(buildUrl(id.toString()));
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

export async function getCoordinatorImages() {
    try {
        const { data } = await axios.get<CoordinatorListEntry[]>(buildUrl("list"));
        return data;
        // return [
        //     {
        //         id: 1,
        //         name: "moin",
        //         imageId: "AF39BBAD222",
        //         signingKeyId: 1,
        //         uniqueId: "ABC123",
        //         signerId: "DEF999"
        //     },
        //     {
        //         id: 2,
        //         name: "cool",
        //         imageId: "AF39BBAD222",
        //         signingKeyId: 1,
        //         uniqueId: "ABC123",
        //         signerId: "DEF999"
        //     },
        //     {
        //         id: 3,
        //         name: "supercoord",
        //         imageId: "AF39BBAD222",
        //         signingKeyId: 1,
        //         uniqueId: "ABC123",
        //         signerId: "DEF999"
        //     }
        // ]
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