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
    this.networkListener = Network.addListener(
      'networkStatusChange',
      (status) => {
        this.haveNetworkConnectivity = status.connected;
      }
    );
  }
}
