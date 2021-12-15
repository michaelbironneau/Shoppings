import { Injectable } from '@angular/core';
import { Network } from '@capacitor/network';
import { PluginListenerHandle } from '@capacitor/core';

@Injectable({
  providedIn: 'root',
})
export class BaseService {
  haveNetworkConnectivity = false;
  networkListener: PluginListenerHandle;
  constructor() {
    Network.getStatus().then((status) => {
      this.haveNetworkConnectivity = status.connected;
    });
    this.networkListener = Network.addListener(
      'networkStatusChange',
      (status) => {
        this.haveNetworkConnectivity = status.connected;
      }
    );
  }
}
