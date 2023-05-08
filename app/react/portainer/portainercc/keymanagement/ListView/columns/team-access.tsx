import { CellProps, Column } from 'react-table';
import { KeyEntry } from '../../types';
import { TeamsSelector } from '@@/TeamsSelector';
import { Form, Formik } from 'formik';
import { updateKey } from '../../keys.service';


export const teamAccess: Column<KeyEntry> = {
  Header: 'Teams',
  accessor: (row) => row.TeamAccessPolicies,
  disableFilters: true,
  Filter: () => null,
  canHide: false,
  disableSortBy: true,
  Cell: AccesCell,
};

export function AccesCell({ value: val, row }: CellProps<KeyEntry>) {
  const initialValues = {
    teamAccessPolicies: Object.keys(val).map(k => Number(k))
  }

  return (
    <>

      <Formik
        initialValues={initialValues}
        onSubmit={() => console.log("x")}
        key={1}
      >
        {({
          values,
          handleSubmit,
          setFieldValue,
        }) => (
          <Form
            className="form-horizontal"
            onSubmit={handleSubmit}
            noValidate
          >
            <TeamsSelector
              value={values.teamAccessPolicies}
              onChange={
                (values) => {
                  setFieldValue('teamAccessPolicies', values)
                  updateAccess(row.original.Id, values)
                }
              }
              teams={row.original.AllTeams}
              placeholder="Select one or more teams to access the key"
            />
          </Form>
        )}
      </Formik>
    </>
  )
}

async function updateAccess(id: any, values: any) {
  let access = values.reduce((prev: any, current: any) => {
    return {
      ...prev,
      [current.toString()]: {
        "RoleId": 0
      }
    }
  }, {})
  console.log(access)
  let data = await updateKey(id, access)
  console.log(data)
}