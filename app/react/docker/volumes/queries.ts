import { useQuery } from "react-query";
import { getEncryptedVolumes } from "./volume.service";

export function useEncryptedVolumes<T = any | undefined>(
  envId: number,
  enabled = true,
  select: (data: any | undefined) => T = (data) => data as unknown as T
) {
  const res = useQuery(
    ['encrypted-volumes'],
    () => getEncryptedVolumes(envId.toString()),
    {
      meta: {
        error: { title: 'Failure', message: 'Unable to load encrypted volumes' },
      },
      enabled,
      select,
    }
  );

  return res;
}