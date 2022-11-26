import { CellProps, Column } from 'react-table';
import { Icon } from '@@/Icon';
import clsx from 'clsx';

import { CoordinatorDeploymentEntry } from '../../types';

export const status: Column<CoordinatorDeploymentEntry> = {
  Header: 'Status',
  accessor: (row) => row.verified,
  disableFilters: true,
  Cell: StatusCell,
  canHide: false,
};

export function StatusCell({ value: status, row }: CellProps<CoordinatorDeploymentEntry>) {
  return (
    <>
      {
        status && <>
          <Icon
            icon="lock"
            className={clsx('icon icon-sm icon-success')}
            feather
          />
          Verified
        </>
      }
      {
        !status &&
        <>
          <Icon
            icon="lock"
            className={clsx('icon icon-sm icon-danger')}
            feather
          />
          Not Verified
        </>
      }
    </>
  );
}