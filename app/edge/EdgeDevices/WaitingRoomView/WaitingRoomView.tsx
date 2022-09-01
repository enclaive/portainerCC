import { useRouter } from '@uirouter/react';

import { useEnvironmentList } from '@/portainer/environments/queries/useEnvironmentList';
import { EdgeTypes } from '@/portainer/environments/types';

import { InformationPanel } from '@@/InformationPanel';
import { TextTip } from '@@/Tip/TextTip';
import { PageHeader } from '@@/PageHeader';

import { Datatable } from './Datatable';

export function WaitingRoomView() {
  const router = useRouter();
  const { environments, isLoading, totalCount } = useEnvironmentList({
    edgeDevice: true,
    edgeDeviceUntrusted: true,
    excludeSnapshots: true,
    types: EdgeTypes,
  });

  if (process.env.PORTAINER_EDITION !== 'BE') {
    router.stateService.go('edge.devices');
    return null;
  }

  return (
    <>
      <PageHeader
        title="Waiting Room"
        breadcrumbs={[
          { label: 'Edge Devices', link: 'edge.devices' },
          { label: 'Waiting Room' },
        ]}
      />

      <InformationPanel>
        <TextTip color="blue">
          Only environments generated from the AEEC script will appear here,
          manually added environments and edge devices will bypass the waiting
          room.
        </TextTip>
      </InformationPanel>

      <Datatable
        devices={environments}
        totalCount={totalCount}
        isLoading={isLoading}
      />
    </>
  );
}
