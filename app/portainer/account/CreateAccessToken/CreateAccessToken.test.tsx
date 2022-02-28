import userEvent from '@testing-library/user-event';

import { renderWithQueryClient } from '@/react-tools/test-utils';
import { server, rest } from '@/setup-tests/server';
import { UserViewModel } from '@/portainer/models/user';
import { UserContext } from '@/portainer/hooks/useUser';

import { CreateAccessToken } from './CreateAccessToken';

test('the button is disabled when description is missing and enabled when description is filled', async () => {
  const queries = renderComponent();

  const button = queries.getByRole('button', { name: 'Add access token' });

  expect(button).toBeDisabled();

  const descriptionField = queries.getByLabelText('Description');

  userEvent.type(descriptionField, 'description');

  expect(button).toBeEnabled();

  userEvent.clear(descriptionField);

  expect(button).toBeDisabled();
});

test('once the button is clicked, the access token is generated and displayed', async () => {
  const token = 'a very long access token that should be displayed';

  const queries = renderComponent(token);

  const descriptionField = queries.getByLabelText('Description');

  userEvent.type(descriptionField, 'description');

  const button = queries.getByRole('button', { name: 'Add access token' });

  userEvent.click(button);

  await expect(queries.findByText('New access token')).resolves.toBeVisible();
  expect(queries.getByText(token)).toHaveTextContent(token);
});

function renderComponent(accessToken = '') {
  server.use(
    rest.post('/api/users/:userId/tokens', (req, res, ctx) =>
      res(ctx.json({ rawAPIKey: accessToken }))
    )
  );

  const user = new UserViewModel({ Username: 'user' });

  const queries = renderWithQueryClient(
    <UserContext.Provider value={{ user }}>
      <CreateAccessToken />
    </UserContext.Provider>
  );

  expect(queries.getByLabelText('Description')).toBeVisible();

  return queries;
}
