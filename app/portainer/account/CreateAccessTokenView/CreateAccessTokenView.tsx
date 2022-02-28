import { PageHeader } from '@/portainer/components/PageHeader';
import { r2a } from '@/react-tools/react2angular';

import { CreateAccessTokenForm } from './CreateAccessTokenForm/CreateAccessTokenForm';

export function CreateAccessTokenView() {
  return (
    <>
      <PageHeader
        title="Create access token"
        breadcrumbs={[
          { label: 'User settings', link: 'portainer.account' },
          { label: 'Add access token' },
        ]}
      />

      <div className="row">
        <div className="col-sm-12">
          <CreateAccessTokenForm />
        </div>
      </div>
    </>
  );
}

export const CreateAccessTokenViewAngular = r2a(CreateAccessTokenView, []);
