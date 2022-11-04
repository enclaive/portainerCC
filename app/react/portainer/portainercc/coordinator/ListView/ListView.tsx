
import { PageHeader } from '@@/PageHeader';
import { Key } from 'react-feather';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';
import { useKeys } from '../../keymanagement/queries';

export function CoordinatorImagesListView() {

    const keysQuery = useKeys('SIGNING')

    let title = "Build your coordinator";


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />

            {keysQuery.data && (
                <BuildCoordinatorForm keys={keysQuery.data} />
            )}

            {/* <Datatable
                columns={columns}
                titleOptions={{
                    title: title,
                    icon: Key,
                }}
                dataset={keysQuery.data}
                settingsStore={store}
                storageKey={storageKey}
                emptyContentLabel="No keys found"
                // isLoading={listQuery.isLoading}
                totalCount={keysQuery.data.length}
                // renderTableActions={(selectedRows) => (
                //     <TableActions selectedRows={selectedRows} />
                // )}
            /> */}
        </>
    );
}
