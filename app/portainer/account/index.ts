import angular from 'angular';

import { CreateAccessTokenViewAngular } from './CreateAccessTokenView';

export const accountModule = angular
  .module('portainer.app.account', [])
  .component('createAccessTokenView', CreateAccessTokenViewAngular).name;
