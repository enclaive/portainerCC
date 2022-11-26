import { useQuery } from "react-query";
import { getKeys } from "./keys.service";

import { KeyEntry } from "./types";

export function useKeys<T = KeyEntry[]>(
    type: string,
    enabled = true,
    select: (data: KeyEntry[]) => T = (data) => data as unknown as T
  ) {
    const keys = useQuery(
      ['keys'],
      () => getKeys(type),
      {
        meta: {
          error: { title: 'Failure', message: 'Unable to load keys' },
        },
        enabled,
        select,
      }
    );
  
    return keys;
  }
  