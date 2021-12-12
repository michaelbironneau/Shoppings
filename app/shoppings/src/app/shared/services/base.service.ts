import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable, OnDestroy } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Network } from '@capacitor/network';
import { from } from 'rxjs';
import { PluginListenerHandle } from '@capacitor/core';

@Injectable({
  providedIn: 'root',
})
export class BaseService implements OnDestroy {
  haveNetworkConnectivity = false;
  networkListener: PluginListenerHandle;
  constructor(protected http: HttpClient) {
    this.networkListener = Network.addListener(
      'networkStatusChange',
      (status) => {
        this.haveNetworkConnectivity = status.connected;
      }
    );
  }

  ngOnDestroy() {
    if (this.networkListener) {
      this.networkListener.remove();
    }
  }
}
