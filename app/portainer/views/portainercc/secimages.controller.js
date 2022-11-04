import angular from 'angular';
import _ from 'lodash-es';

angular.module('portainer.app').controller('secImagesController', secImagesController);

/* @ngInject */
export default function secImagesController(Notifications, $q, $scope, SecImagesService) {

  $scope.state = {
    actionInProgress: false,
  };

  function initView() {

    $q.all({
      images: SecImagesService.getImages(),
    })
      .then(function success(data) {
        $scope.imageIdentifiers = _.orderBy(data.images, 'image', 'asc');
        console.log("MOIN");
        console.log($scope.data)
      }).catch(function error(err) {
        $scope.imageIdentifiers = [];
        Notifications.error('Failure', err, 'Unable to retrieve keys');
      })

  }


  initView();
}
