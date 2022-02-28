import angular from 'angular';

import { CreateAccessTokenAngular } from './CreateAccessToken';

export const accountModule = angular
  .module('portainer.app.account', [])
  .component('createAccessToken', CreateAccessTokenAngular).name;
