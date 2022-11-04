angular.module('portainer.app').factory('Secureimages', [
  '$resource',
  'API_ENDPOINT_SECIMAGES',
  function SecImagesFactory($resource, API_ENDPOINT_SECIMAGES) {
    'use strict';
    return $resource(
      API_ENDPOINT_SECIMAGES,
      {},
      {
        query: {
          method: 'GET', isArray: true
        },
      }
    );
  },
]);
