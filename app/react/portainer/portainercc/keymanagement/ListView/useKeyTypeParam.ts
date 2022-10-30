import { useCurrentStateAndParams } from '@uirouter/react';

export function useKeyTypeParam() {
  const {
    params: { type: keytype },
  } = useCurrentStateAndParams();

  return keytype;
}
