import { useQuery } from "react-query";
import { KeyId } from "../keymanagement/types";
import { buildCoordinator } from "./coordinator.service";


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
  