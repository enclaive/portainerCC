
import { PageHeader } from '@@/PageHeader';
import { Key, Download, Trash2 } from 'react-feather';

import { CreateKeyForm } from './CreateKeyForm/CreateKeyForm';
import { Datatable } from '@@/datatables';
import { columns } from './columns';
import { Button } from '@@/buttons';
import { createStore } from '@/react/portainer/environments/update-schedules/ListView/datatable-store';

import { useKeyTypeParam } from './useKeyTypeParam';
import { useTeams } from '@/react/portainer/users/teams/queries';
import { useUsers } from '@/portainer/users/queries';
import { useUser } from '@/portainer/hooks/useUser';
import { useKeys } from '../queries';
import { KeyEntry } from '../types';

const storageKey = 'portainercc-keys';
const useStore = createStore(storageKey);

export function KeyListView() {
    const { isAdmin } = useUser();

    const store = useStore();

    const type = useKeyTypeParam();
    const usersQuery = useUsers(false);
    const teamsQuery = useTeams(!isAdmin, 0, { enabled: !!usersQuery.data });
    const keysQuery = useKeys(type)
    let title = "";

    if (type == "SIGNING") {
        title = "SGX Signing Keys"
    } else if (type == "FILE_ENC") {
        title = "Gramine Protected File Keys"
    } else {
        throw Error("invalid key type")
    }

    if (!keysQuery.data) {
        return null;
    }
    console.log(keysQuery)

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {teamsQuery.data && (
                <CreateKeyForm keytype={type} teams={teamsQuery.data} />
            )}

            <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Key,
                }}
                dataset={keysQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                isLoading={keysQuery.isLoading}
                totalCount={keysQuery.data.length}
                renderTableActions={(selectedRows) => (
                    <TableActions selectedRows={selectedRows} />
                )}
            />
        </>
    );
}

function TableActions({ selectedRows }: { selectedRows: KeyEntry[] }) {
    return (
        <div>
            <Button
                icon={Download}
                color="primary"
                disabled={selectedRows.length === 0}
                onClick={() => handleExport()}
            >
                Export
            </Button>
            <Button
                icon={Trash2}
                color="dangerlight"
                disabled={selectedRows.length === 0}
                onClick={() => handleRemove()}
            >
                Remove
            </Button>
        </div>
    );

    function handleRemove() {
        const ids = selectedRows.map((row) => row.Id);
        console.log("REMOVE:")
        console.log(ids)
    }

    function handleExport() {
        const ids = selectedRows.map((row) => row.Id);
        console.log("EXPORT:")
        console.log(ids)
    }
}