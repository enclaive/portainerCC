import angular from 'angular';

import { CreateAccessTokenAngular } from './CreateAccessTokenForm';

export const accountModule = angular
  .module('portainer.app.account', [])
  .component('createAccessToken', CreateAccessTokenAngular).name;
