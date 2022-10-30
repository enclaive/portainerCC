import { useQuery } from "react-query";


// export function useKeys<T = KeyEntry[]>(
//     type: string,
//     enabled = true,
//     select: (data: KeyEntry[]) => T = (data) => data as unknown as T
//   ) {
//     // const keys = useQuery(
//     //   ['keys'],
//     //   () => getKeys(type),
//     //   {
//     //     meta: {
//     //       error: { title: 'Failure', message: 'Unable to load keys' },
//     //     },
//     //     enabled,
//     //     select,
//     //   }
//     // );
  
//     // return keys;

//     return null;
//   }
  