import angular from 'angular';
import _ from 'lodash-es';

angular.module('portainer.app').controller('keylistController', keylistController);

/* @ngInject */
export default function keylistController(Notifications, $q, $scope, KeymanagementService, TeamService, $state, FileSaver, $stateParams) {

  $scope.state = {
    actionInProgress: false,
  };

  var KEY_TYPE = "";

  $scope.keyTitle = "";

  var tempTeamIds = [];

  $scope.formData = {
    description: "",
    teamIds: [],
  }

  this.generateKey = function () {
    $scope.state.actionInProgress = true;

    var teamIds = $scope.formData.teamIds.map((team) => { return team.Id });

    KeymanagementService.createKey(KEY_TYPE, $scope.formData.description, teamIds, null)
      .then(function success() {
        Notifications.success('Success', 'New Key added!');
        $state.reload();
      }).catch(function error(err) {
        Notifications.error('Failure', err, 'Unable to generate key');
      })
      .finally(function final() {
        $scope.state.actionInProgress = false;
      });
  }

  this.removeKey = function (selectedKeys) {
    $scope.state.actionInProgress = true;

    $q.all(
      selectedKeys.map(async key => {
        await KeymanagementService.deleteKey(key.id)
          .then(function success() {
            Notifications.success('Success', 'Key deleted!');
          })
          .catch(function error(err) {
            Notifications.error('Failure', err, 'Unable to delete key!');
          })
      })
    )
      .then(function success() {
        $scope.state.actionInProgress = false;
        $state.reload();
      })
  }

  this.exportKey = function (selectedKeys) {
    KeymanagementService.getKeyAsPEM(selectedKeys[0].id)
      .then(function success(data) {
        console.log(data);
        var downloadData = new Blob([data.PEM], { type: 'text/plain' });
        FileSaver.saveAs(downloadData, 'enclave_signing_key_' + data.Id + '.pem');
        Notifications.success('Key successfully exported');
      })
      .catch(function error(err) {
        Notifications.error('Failure', err, 'Unable to export key');
      })

  }


  this.importKey = function (file) {
    $scope.state.actionInProgress = true;
    readFileContent(file).then(function success(pem) {
      var teamIds = $scope.formData.teamIds.map((team) => { return team.Id });

      KeymanagementService.createKey(KEY_TYPE, $scope.formData.description, teamIds, pem)
        .then(function success() {
          Notifications.success('Success', 'New Key imported!');
          $state.reload();
        }).catch(function error(err) {
          Notifications.error('Failure', err, 'Unable to import key');
        })
        .finally(function final() {
          $scope.state.actionInProgress = false;
        });
    })

  }

  this.updateKeyAccess = function (key) {
    var newTeamIds = key.teamsSelection.map((team) => { return team.Id })
    if (!_.isEqual(tempTeamIds, newTeamIds)) {
      $scope.state.actionInProgress = true;
      KeymanagementService.updateTeams(key.id, newTeamIds)
        .then(function success() {
          Notifications.success('Success', 'Access updated!');
        })
        .catch(function error(err) {
          Notifications.error('Failure', err, 'Unable to update access!');
        })
        .finally(function final() {
          $scope.state.actionInProgress = false;
        });
    }
    tempTeamIds = [];
  }

  this.saveTempSelection = function (key) {
    tempTeamIds = key.teamsSelection.map((team) => { return team.Id })
  }


  function readFileContent(file) {
    return new Promise((resolve, reject) => {
      var fr = new FileReader();
      fr.onload = () => {
        resolve(fr.result);
      }
      fr.onerror = reject;
      fr.readAsText(file);
    })
  }

  function initView() {
    $q.all({
      keys: KeymanagementService.getKeys(KEY_TYPE),
      teams: TeamService.teams()
    })
      .then(function success(data) {
        var keys = _.orderBy(data.keys, 'description', 'asc');

        $scope.keys = keys.map((key) => {
          key.teams = angular.copy(data.teams)

          if (!_.isEmpty(key.TeamAccessPolicies)) {
            key.teams = key.teams.map((team) => {
              if (Object.keys(key.TeamAccessPolicies).includes(team.Id.toString())) {
                team.ticked = true;
              }
              return team;
            })
          }
          return key
        })

        $scope.teams = _.orderBy(data.teams, 'Name', 'asc');
      }).catch(function error(err) {
        $scope.keys = [];
        $scope.teams = [];
        Notifications.error('Failure', err, 'Unable to retrieve keys');
      })

  }


  KEY_TYPE = $stateParams.type;
  if (KEY_TYPE == "SIGNING") {
    $scope.keyTitle = "Enclave Signing"
    $scope.keySubtitle = "Manage your Signing Keys to build SGX enhanced containers"
    initView();
  }
  else if (KEY_TYPE == "FILE_ENC") {
    $scope.keyTitle = "File Encryption"
    $scope.keySubtitle = "Manage your Keys to encrypt files"
    initView();
  } else {
    $state.go('portainer.home', {}, { reload: true });
  }
}
