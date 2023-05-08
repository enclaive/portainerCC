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

export async function getKey(id: number) {
  try {
    const { data } = await axios.get<KeyEntry>(buildUrl(id));
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function deleteKey(id: number) {
  try {
    const { data } = await axios.delete<KeyEntry>(buildUrl(id));
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function updateKey(id: number, access: any) {
  try {
    let payload: any = {
      TeamAccessPolicies: access,
    };

    const { data } = await axios.post<KeyEntry[]>(buildUrl(id), payload);
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function createKey(
  type: string,
  desc: string,
  access: any,
  file?: string
) {
  try {
    let payload: any = {
      KeyType: type,
      Description: desc,
      TeamAccessPolicies: access,
    };

    if (file) {
      payload.Data = file;
    }

    const { data } = await axios.post<KeyEntry[]>(buildUrl(), payload);
    return data;
  } catch (error) {
    throw parseAxiosError(error as Error);
  }
}

export async function importKey(
  type: string,
  desc: string,
  access: any,
  data: string
) {
  return createKey(type, desc, access, data);
}

function buildUrl(id?: KeyId) {
  let url = '/portainercc/keys';

  if (id) {
    url += `/${id}`;
  }

  return url;
}
