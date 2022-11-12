import { useQuery } from "react-query";
import { KeyId } from "../keymanagement/types";
import { buildCoordinator, getCoordinatorImages } from "./coordinator.service";
import { CoordinatorListEntry } from "./types";


export function useBuildCoordinator(name: string, keyId: KeyId) {
    return useQuery(
      ['coordinator'],
      () => buildCoordinator(name, keyId),
      {
        meta: {
          error: { title: 'Failure', message: 'Unable to build coordinator' },
        },
      }
    );
  }
  
  export function useCoordinatorImages<T = CoordinatorListEntry[]>(
    enabled = true,
    select: (data: CoordinatorListEntry[]) => T = (data) => data as unknown as T
  ) {
    const images = useQuery(
      ['coordinator-images'],
      () => getCoordinatorImages(),
      {
        meta: {
          error: { title: 'Failure', message: 'Unable to load coordinator images' },
        },
        enabled,
        select,
      }
    );
  
    return images;
  }
  