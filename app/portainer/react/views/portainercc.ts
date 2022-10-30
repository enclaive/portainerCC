import angular from 'angular';
import { StateRegistry } from '@uirouter/angularjs';

import { KeyListView } from '@/react/portainer/portainercc/keymanagement';
import { CoordinatorImagesListView } from '@/react/portainer/portainercc/coordinator';

import { r2a } from '@/react-tools/react2angular';
import { withCurrentUser } from '@/react-tools/withCurrentUser';
import { withReactQuery } from '@/react-tools/withReactQuery';
import { withUIRouter } from '@/react-tools/withUIRouter';

export const portainerCCModule = angular
  .module('portainer.app.portainercc', [])
  .config(config)
  .component(
    'raList',
    r2a(withUIRouter(withReactQuery(withCurrentUser(KeyListView))), [])
  )
  .component(
    'keymanagement',
    r2a(withUIRouter(withReactQuery(withCurrentUser(KeyListView))), [])
  ).component(
    'coordinatorBuild',
    r2a(withUIRouter(withReactQuery(withCurrentUser(CoordinatorImagesListView))), [])
  ).name;

/* @ngInject */
function config($stateRegistryProvider: StateRegistry) {
  $stateRegistryProvider.register({
    name: 'portainer.keymanagement',
    url: '/keys?type',
    views: {
      'content@': {
        component: 'keymanagement',
      },
    },
  });

  $stateRegistryProvider.register({
    name: 'portainer.raList',
    url: '/remote-attestation',
    views: {
      'content@': {
        component: 'raList',
      },
    },
  });

  $stateRegistryProvider.register({
    name: 'portainer.coordinator-build',
    url: '/coordinator',
    views: {
      'content@': {
        component: 'coordinatorBuild',
      },
    },
  });
}
