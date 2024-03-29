import _ from 'lodash-es';
import { getKeys } from '@/react/portainer/portainercc/keymanagement/keys.service';
import { AccessControlFormData } from '../../../../portainer/components/accessControlForm/porAccessControlFormModel';
import { VolumesNFSFormData } from '../../../components/volumesNFSForm/volumesNFSFormModel';
import { VolumesCIFSFormData } from '../../../components/volumesCIFSForm/volumesCifsFormModel';


angular.module('portainer.docker').controller('CreateVolumeController', [
  '$q',
  '$scope',
  '$state',
  'VolumeService',
  'PluginService',
  'ResourceControlService',
  'Authentication',
  'Notifications',
  'FormValidator',
  'HttpRequestHelper',
  function ($q, $scope, $state, VolumeService, PluginService, ResourceControlService, Authentication, Notifications, FormValidator, HttpRequestHelper) {
    $scope.formValues = {
      Driver: 'local',
      DriverOptions: [],
      AccessControlData: new AccessControlFormData(),
      NodeName: null,
      NFSData: new VolumesNFSFormData(),
      CIFSData: new VolumesCIFSFormData(),
      usePF: false,
      selectedPFKey: null,
    };

    $scope.state = {
      formValidationError: '',
      actionInProgress: false,
    };

    $scope.availableVolumeDrivers = [];

    $scope.addDriverOption = function () {
      $scope.formValues.DriverOptions.push({ name: '', value: '' });
    };

    $scope.removeDriverOption = function (index) {
      $scope.formValues.DriverOptions.splice(index, 1);
    };

    $scope.onUseNFSChange = onUseNFSChange;

    function onUseNFSChange(checked) {
      return $scope.$evalAsync(() => {
        $scope.formValues.NFSData.useNFS = checked;
        $scope.formValues.CIFSData.useCIFS = false;
      });
    }

    $scope.onUseCIFSChange = onUseCIFSChange;

    function onUseCIFSChange(checked) {
      return $scope.$evalAsync(() => {
        $scope.formValues.CIFSData.useCIFS = checked;
        $scope.formValues.NFSData.useNFS = false;
      });
    }

    function validateForm(accessControlData, isAdmin) {
      $scope.state.formValidationError = '';
      var error = '';
      error = FormValidator.validateAccessControl(accessControlData, isAdmin);

      if (error) {
        $scope.state.formValidationError = error;
        return false;
      }
      return true;
    }

    function prepareCIFSConfiguration(driverOptions) {
      const data = $scope.formValues.CIFSData;

      driverOptions.push({ name: 'type', value: 'cifs' });

      let share = data.share.replace('\\', '/');
      if (share[0] !== '/') {
        share = '/' + share;
      }
      const device = '//' + data.serverAddress + share;
      driverOptions.push({ name: 'device', value: device });

      const versionNumber = data.versionsNumber[data.versions.indexOf(data.version)];
      const options = 'addr=' + data.serverAddress + ',username=' + data.username + ',password=' + data.password + ',vers=' + versionNumber;
      driverOptions.push({ name: 'o', value: options });
    }

    function prepareNFSConfiguration(driverOptions) {
      var data = $scope.formValues.NFSData;

      driverOptions.push({ name: 'type', value: 'nfs' });

      var options = 'addr=' + data.serverAddress + ',' + data.options;
      if (data.version === 'NFS4') {
        options = options + ',nfsvers=4';
      }
      driverOptions.push({ name: 'o', value: options });

      var mountPoint = data.mountPoint[0] === ':' ? data.mountPoint : ':' + data.mountPoint;
      driverOptions.push({ name: 'device', value: mountPoint });
    }

    $scope.create = function () {
      var name = $scope.formValues.Name;
      var driver = $scope.formValues.Driver;
      var driverOptions = $scope.formValues.DriverOptions;

      if ($scope.formValues.NFSData.useNFS) {
        prepareNFSConfiguration(driverOptions);
      }

      if ($scope.formValues.CIFSData.useCIFS) {
        prepareCIFSConfiguration(driverOptions);
      }

      var volumeConfiguration = VolumeService.createVolumeConfiguration(name, driver, driverOptions);
      var accessControlData = $scope.formValues.AccessControlData;
      var userDetails = Authentication.getUserDetails();
      var isAdmin = Authentication.isAdmin();

      if (!validateForm(accessControlData, isAdmin)) {
        return;
      }

      var nodeName = $scope.formValues.NodeName;
      HttpRequestHelper.setPortainerAgentTargetHeader(nodeName);

      $scope.state.actionInProgress = true;

      //if gramine encryption enabled, add encryption and key id as volume labels
      if ($scope.formValues.usePF && $scope.formValues.selectedPFKey) {
        //create new key or use existing
        if ($scope.formValues.selectedPFKey.Id == 0) {
          volumeConfiguration.createKey = true;
        } else {
          volumeConfiguration.Labels = {
            encrypted: "true",
            pfEncryptionKeyId: $scope.formValues.selectedPFKey.Id.toString()
          }
        }
      }

      VolumeService.createVolume(volumeConfiguration)
        .then(function success(data) {
          const userId = userDetails.ID;
          const resourceControl = data.ResourceControl;
          return ResourceControlService.applyResourceControl(userId, accessControlData, resourceControl);
        })
        .then(function success() {
          Notifications.success('Success', 'Volume successfully created');
          $state.go('docker.volumes', {}, { reload: true });
        })
        .catch(function error(err) {
          Notifications.error('Failure', err, 'An error occured during volume creation');
        })
        .finally(function final() {
          $scope.state.actionInProgress = false;
        });
    };

    function initView() {
      var apiVersion = $scope.applicationState.endpoint.apiVersion;

      $q.all({
        keys: getKeys("FILE_ENC"),
        volumes: PluginService.volumePlugins(apiVersion < 1.25)
      })
        .then(function success(data) {
          $scope.pfKeys = [{ Id: 0, Description: "Create new key" }, ..._.orderBy(data.keys, 'Description', 'asc').map((key) => { return { Id: key.Id, Description: key.Description } })];
          console.log($scope.pfKeys)
          $scope.availableVolumeDrivers = data.volumes;

        }).catch(function error(err) {
          Notifications.error('Failure', err, 'Unable to retrieve keys or volume drivers');
        })
    }

    initView();
  },
]);
