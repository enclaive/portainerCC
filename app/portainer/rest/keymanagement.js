angular.module('portainer.app').factory('Keymanagement', [
  '$resource',
  'API_ENDPOINT_KEYMANAGEMENT',
  function KeymanagementFactory($resource, API_ENDPOINT_KEYMANAGEMENT) {
    'use strict';
    return $resource(
      API_ENDPOINT_KEYMANAGEMENT + '/:id/:action',
      {},
      {
        create: { method: 'POST', ignoreLoadingBar: true },
        query: { method: 'GET', params: { type: '@type' }, isArray: true },
        export: { method: 'GET', params: { id: '@id' } },
        update: { method: 'PUT', params: { id: '@id' } },
        delete: { method: 'DELETE', params: { id: '@id' } },
      }
    );
  },
]);
