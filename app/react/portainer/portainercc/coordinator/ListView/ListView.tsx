
import { PageHeader } from '@@/PageHeader';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';
import { useKeys } from '../../keymanagement/queries';
import { Datatable } from '@@/datatables';
import { columns } from './columns';
import { deploymentColumns } from './deployment-columns';
import { Codesandbox, Trash2 } from 'react-feather';
import { createStore } from '@/react/portainer/environments/update-schedules/ListView/datatable-store';
import { useCoordinatorImages } from '../queries';
import { Button } from '@@/buttons';
import { CoordinatorListEntry } from '../types';
import { removeCoordinatorImage } from '../coordinator.service';
import { useCoordinatorDeployments } from '@/react/docker/portainercc/coordinator/queries';
import { useEnvironmentList } from '@/portainer/environments/queries/useEnvironmentList';

const storageKey = 'portainercc-coordinators';
const useStore = createStore(storageKey);

export function CoordinatorImagesListView() {

    const store = useStore();

    const keysQuery = useKeys('SIGNING')
    const coordintaorQuery = useCoordinatorImages();
    const deploymentQuery = useCoordinatorDeployments()
    const envQuery = useEnvironmentList()

    let title = "Coordinator";

    if (!coordintaorQuery.data) {
        return null;
    }

    if (!envQuery.environments) {
        return null;
    }

    if (!deploymentQuery.data) {
        return null;
    }

    deploymentQuery.data = deploymentQuery.data.map(entry => {
        entry.endpointName = envQuery.environments.find(e => e.Id == entry.endpointId)?.Name
        return entry
    })

    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {keysQuery.data && (
                <BuildCoordinatorForm keys={keysQuery.data} />
            )}


            <Datatable
                columns={columns}
                titleOptions={{
                    title: "Coordinator Images",
                    icon: Codesandbox,
                }}
                dataset={coordintaorQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No coordinator images found"
                isLoading={coordintaorQuery.isLoading}
                totalCount={coordintaorQuery.data.length}
                renderTableActions={(selectedRows) => (
                    <TableActions selectedRows={selectedRows} />
                )}
            />

            <Datatable
                columns={deploymentColumns}
                titleOptions={{
                    title: "Deployed Coordinators",
                    icon: Codesandbox,
                }}
                disableSelect
                dataset={deploymentQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                isLoading={deploymentQuery.isLoading}
                totalCount={deploymentQuery.data.length}
            />

        </>
    );
}

function TableActions({ selectedRows }: { selectedRows: CoordinatorListEntry[] }) {
    return (
        <div>
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
        for await (const row of selectedRows) {
            await removeCoordinatorImage(row.id)
        }
        return;
    }
}