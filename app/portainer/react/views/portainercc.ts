import angular from 'angular';
import { StateRegistry } from '@uirouter/angularjs';

import { KeyListView } from '@/react/portainer/portainercc/keymanagement';
import { ConfidentialImagesListView } from '@/react/portainer/portainercc/confidential-images';
import { CoordinatorImagesListView } from '@/react/portainer/portainercc/coordinator';
import { CoordinatorDeploymentView } from '@/react/docker/portainercc/coordinator/DeploymentView';
import { ConfidentialTemplatesView } from '@/react/docker/portainercc/confidential-templates/DeploymentView';
import { RunYourCodeView } from '@/react/docker/portainercc/runyourcode/RunYourCodeView';

import { r2a } from '@/react-tools/react2angular';
import { withCurrentUser } from '@/react-tools/withCurrentUser';
import { withReactQuery } from '@/react-tools/withReactQuery';
import { withUIRouter } from '@/react-tools/withUIRouter';

export const portainerCCModule = angular
  .module('portainer.app.portainercc', [])
  .config(config)
  .component(
    'confimages',
    r2a(withUIRouter(withReactQuery(withCurrentUser(ConfidentialImagesListView))), [])
  )
  .component(
    'keymanagement',
    r2a(withUIRouter(withReactQuery(withCurrentUser(KeyListView))), [])
  ).component(
    'coordinatorBuild',
    r2a(withUIRouter(withReactQuery(withCurrentUser(CoordinatorImagesListView))), [])
  ).component(
    'coordinatorDeployment',
    r2a(withUIRouter(withReactQuery(withCurrentUser(CoordinatorDeploymentView))), [])
  ).component(
    'confidentialTemplates',
    r2a(withUIRouter(withReactQuery(withCurrentUser(ConfidentialTemplatesView))), [])
  ).component(
    'runyourcode',
    r2a(withUIRouter(withReactQuery(withCurrentUser(RunYourCodeView))), [])
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
    name: 'portainer.confimages',
    url: '/confidential-images',
    views: {
      'content@': {
        component: 'confimages',
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

  $stateRegistryProvider.register({
    name: 'docker.coordinator',
    url: '/coordinator',
    views: {
      'content@': {
        component: 'coordinatorDeployment',
      },
    },
  });

  $stateRegistryProvider.register({
    name: 'docker.templates.confidential',
    url: '/confidential',
    views: {
      'content@': {
        component: 'confidentialTemplates',
      },
    },
  });

  $stateRegistryProvider.register({
    name: 'docker.runyourcode',
    url: '/runyourcode',
    views: {
      'content@': {
        component: 'runyourcode',
      },
    },
  });
}
