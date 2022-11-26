import axios, { parseAxiosError } from '@/portainer/services/axios';
import { KeyEntry, KeyId } from './types';

export async function getKeys(type: string) {
  try {
    const { data } = await axios.get<KeyEntry[]>(buildUrl(), {
      params: { type },
    });
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function createKey(type: string, desc: string, access: any, hex?: string) {
  try {
    let payload: any = {
      KeyType: type,
      Description: desc,
      TeamAccessPolicies: access
    }

    if (hex) {
      payload.Data = hex
    }

    const { data } = await axios.post<KeyEntry[]>(buildUrl(), payload);
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function importKey(type: string, desc: string, access: any, hex: string) {
  return createKey(type, desc, access, hex)
}

function buildUrl(id?: KeyId) {
  let url = '/portainercc/keys';

  if (id) {
    url += `/${id}`;
  }

  return url;
}