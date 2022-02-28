import { react2angular } from '@/react-tools/react2angular';

import { CreateAccessTokenForm } from './CreateAccessTokenForm';

const CreateAccessTokenAngular = react2angular(CreateAccessTokenForm, []);

export { CreateAccessTokenForm as CreateAccessToken, CreateAccessTokenAngular };
