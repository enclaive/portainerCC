
angular.module('portainer.app').factory('KeymanagementService', [
  '$q',
  'Keymanagement',
  function KeymanagementServiceFactory($q, Keymanagement) {
    'use strict';
    var service = {};

    service.createKey = function (type, description, teamIds, data) {

      var deferred = $q.defer();

      var teamAccessPolicies = teamIds.reduce((acc, id) => ({ ...acc, [id]: { RoleId: 0 } }), {})

      var payload = {
        KeyType: type,
        Description: description,
        TeamAccessPolicies: teamAccessPolicies
      }

      if (data) {
        payload.data = data;
      }

      Keymanagement.create({}, payload)
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to generate key', err: err })
        })

      return deferred.promise;
    }

    service.updateTeams = function (id, teamIds) {
      var deferred = $q.defer();

      var teamAccessPolicies = teamIds.reduce((acc, id) => ({ ...acc, [id]: { RoleId: 0 } }), {})

      var payload = {
        TeamAccessPolicies: teamAccessPolicies,
      };

      Keymanagement.update({ id: id }, payload)
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to update key', err: err })
        })

      return deferred.promise;
    };

    service.deleteKey = function (keyId) {
      var deferred = $q.defer();

      Keymanagement.delete({ id: keyId })
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to delete key', err: err })
        })

      return deferred.promise;
    }

    service.exportKey = function (keyId) {
      var deferred = $q.defer();
      Keymanagement.export({ id: keyId })
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to retrieve key', err: err })
        });
      return deferred.promise;
    }

    service.getKeys = function (type) {
      var deferred = $q.defer();
      Keymanagement.query({ type: type })
        .$promise.then(function success(data) {
          deferred.resolve(data);
        }).catch(function error(err) {
          deferred.reject({ msg: 'Unable to retrieve keys', err: err })
        });
      return deferred.promise;
    }

    return service;
  },
]);
