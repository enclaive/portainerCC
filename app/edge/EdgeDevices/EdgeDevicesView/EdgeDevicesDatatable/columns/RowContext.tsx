import { EnvironmentGroup } from '@/portainer/environment-groups/types';

import { createRowContext } from '@@/datatables/RowContext';

interface RowContextState {
  isOpenAmtEnabled: boolean;
  groups: EnvironmentGroup[];
}

const { RowProvider, useRowContext } = createRowContext<RowContextState>();

export { RowProvider, useRowContext };
