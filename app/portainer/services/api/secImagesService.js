
angular.module('portainer.app').factory('SecImagesService', [
  '$q',
  'Secureimages',
  function SecImagesServiceFactory($q, Secureimages) {
    'use strict';
    var service = {};

    service.getImages = function () {
      var deferred = $q.defer();
      Secureimages.query()
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to retrieve image list', err: err })
        });
      return deferred.promise;
    }

    return service;
  },
]);
