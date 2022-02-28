import { useMutation } from 'react-query';

import axios, { parseAxiosError } from '@/portainer/services/axios';
import { UserId } from '@/portainer/users/types';
import { error as notifyError } from '@/portainer/services/notifications';

interface CreateAccessTokenResponse {
  rawAPIKey: string;
}

export async function createAccessToken(id: UserId, description: string) {
  try {
    const { data } = await axios.post<CreateAccessTokenResponse>(
      `/users/${id}/tokens`,
      { description }
    );
    return data.rawAPIKey;
  } catch (error) {
    throw parseAxiosError(error as Error, 'Unable to create access token');
  }
}

export function useCreateAccessTokenMutation() {
  return useMutation(
    ({ id, description }: { id: UserId; description: string }) =>
      createAccessToken(id, description),
    {
      onError(error) {
        notifyError('Failure', error as Error, 'Unable to create access token');
      },
    }
  );
}
