import axios, { parseAxiosError } from '@/portainer/services/axios';
import { KeyId } from '../keymanagement/types';

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
        console.log("hier im axioas")
        const { data } = await axios.post(buildUrl(undefined, "build"), { Name: name, SigningKeyId: keyId });
        return data;
    } catch (e) {
        throw parseAxiosError(e as Error, 'Unable to build coordinator')
    }
}

function buildUrl(id?: KeyId, action?: string) {
    let url = '/ra/coordinator';

    if (id) {
        url += `/${id}`;
    }

    if (action) {
        url += `/${action}`;
    }

    return url;
}