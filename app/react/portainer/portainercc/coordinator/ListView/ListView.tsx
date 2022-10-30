
import { PageHeader } from '@@/PageHeader';
import { Key } from 'react-feather';
import { BuildCoordinatorForm } from './BuildCoordinatorForm/BuildCoordinatorForm';


export function CoordinatorImagesListView() {


    let title = "Build your coordinator";


    return (
        <>
            <PageHeader title={title} breadcrumbs={[{ label: 'PortainerCC' }]} />
                <BuildCoordinatorForm />

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
