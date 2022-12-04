import { CellProps, Column } from 'react-table';


import { CoordinatorDeploymentEntry } from '../../types';

export const environment: Column<CoordinatorDeploymentEntry> = {
  Header: 'Environment',
  accessor: (row) => row,
  Cell: EnvCell,
  disableFilters: true,
  canHide: false,
};

export function EnvCell({ value: row }: CellProps<CoordinatorDeploymentEntry>) {
  return (
    <>
      {row.endpointName}
    </>
  );
}