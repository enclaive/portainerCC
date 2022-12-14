import _ from 'lodash-es';

export default class KubernetesApplicationIngressController {
  /* @ngInject */
  constructor($async, KubernetesIngressService) {
    this.$async = $async;
    this.KubernetesIngressService = KubernetesIngressService;
  }

  $onInit() {
    return this.$async(async () => {
      this.hasIngress;
      this.applicationIngress = [];
      const ingresses = await this.KubernetesIngressService.get(this.application.ResourcePool);
      const services = this.application.Services;

      _.forEach(services, (service) => {
        _.forEach(ingresses, (ingress) => {
          _.forEach(ingress.Paths, (path) => {
            if (path.ServiceName === service.metadata.name) {
              path.Secure = ingress.TLS && ingress.TLS.filter((tls) => tls.hosts && tls.hosts.includes(path.Host)).length > 0;
              this.applicationIngress.push(path);
              this.hasIngress = true;
            }
          });
        });
      });
    });
  }
}
