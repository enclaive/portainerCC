
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
import { deleteKey, getKey } from '../keys.service';
import FileSaver from 'file-saver';

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

    if (!keysQuery.data || !teamsQuery.data) {
        return null;
    }

    let teams = teamsQuery.data

    //TODO stupid workaround
    let entries = keysQuery.data.map(e => {
        e.AllTeams = teams
        return e
    })

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            <CreateKeyForm keytype={type} teams={teams} />

            <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Key,
                }}
                dataset={entries}
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
                disabled={selectedRows.length !== 1}
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

    async function handleRemove() {
        selectedRows.map(async (row) => {
            const data = await deleteKey(row.Id)
            console.log(data)
        });
    }

    async function handleExport() {
        //its only 1
        const id = selectedRows.map((row) => row.Id)[0];
        const data = await getKey(id);
        console.log(data)
        if (data && data.Export) {
            const fileContent = new Blob([data.Export], { type: 'text/plain' });
            let fileName = data.Description
            if(data.KeyType == "FILE_ENC"){
                fileName = fileName + ".pfkey"
            }else {
                fileName = fileName + ".pem"
            }
            FileSaver.saveAs(fileContent, fileName)
        }
    }
}