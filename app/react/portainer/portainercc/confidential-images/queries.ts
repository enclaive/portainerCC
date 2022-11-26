import { useQuery } from "react-query";
import { getConfImages } from "./confimage.service";
import { ConfidentialImage } from "./types";


  export function useConfImages<T = ConfidentialImage[]>(
    enabled = true,
    select: (data: ConfidentialImage[]) => T = (data) => data as unknown as T
  ) {
    const images = useQuery(
      ['conf-images'],
      () => getConfImages(),
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
  