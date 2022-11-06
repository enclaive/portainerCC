
import { PageHeader } from '@@/PageHeader';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';
import { useKeys } from '../../keymanagement/queries';
import { Datatable } from '@@/datatables';
import { columns } from './columns';
import { Codesandbox } from 'react-feather';
import { CoordinatorListEntry } from '../types'
import { createStore } from '@/react/portainer/environments/update-schedules/ListView/datatable-store';

const storageKey = 'portainercc-coordinators';
const useStore = createStore(storageKey);

export function CoordinatorImagesListView() {

    const store = useStore();

    const keysQuery = useKeys('SIGNING')
    const coordintaorQuery = null;

    const exampleCoordinatorResult: CoordinatorListEntry[] = [
        {
            id: 1,
            name: "moin",
            imageId: "AF39BBAD222",
            signingKeyId: 1,
            uniqueId: "ABC123",
            signerId: "DEF999"
        },
        {
            id: 2,
            name: "cool",
            imageId: "AF39BBAD222",
            signingKeyId: 1,
            uniqueId: "ABC123",
            signerId: "DEF999"
        }
    ]

    let title = "Coordinator images";


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {keysQuery.data && (
                <BuildCoordinatorForm keys={keysQuery.data} />
            )}

            <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Codesandbox,
                }}
                disableSelect
                dataset={exampleCoordinatorResult}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                // isLoading={listQuery.isLoading}
                totalCount={exampleCoordinatorResult.length}
                // renderTableActions={(selectedRows) => (
                //     <TableActions selectedRows={selectedRows} />
                // )}
            />
        </>
    );
}
