import { HttpClient } from '@angular/common/http';
import { Injectable, OnDestroy } from '@angular/core';
import { Network } from '@capacitor/network';
import { PluginListenerHandle } from '@capacitor/core';
import { ListItem } from '../models/list-item';
import { environment } from 'src/environments/environment';
import {
  ListUpdate,
  PriorityQueueItem,
  PRIORITY_UPDATE_ARCHIVE,
  PRIOTITY_UPDATE_NEW,
} from '../models/list-update';
import { BehaviorSubject, forkJoin } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class BaseService implements OnDestroy {
  haveNetworkConnectivity = false;
  networkListener: PluginListenerHandle;
  syncing: BehaviorSubject<boolean>;
  constructor(protected http: HttpClient) {
    this.syncing.next(false);
    this.networkListener = Network.addListener(
      'networkStatusChange',
      (status) => {
        if (status.connected && !this.haveNetworkConnectivity) {
          this.pushQueuedUpdates();
        }
        this.haveNetworkConnectivity = status.connected;
      }
    );
  }

  pushQueuedUpdates() {
    if (!environment.api) {
      console.warn(
        'Cannot push queued updates, not sure why this has been invoked!'
      );
      localStorage.removeItem('queue'); // remove this for good measure if we're in demo mode
    }
    if (!this.haveNetworkConnectivity) {
      console.warn('Cannot push queued updates without network connectivity');
      return;
    }
    this.syncing.next(true);
    const requests$ = [];
    const queue = localStorage.getItem('queue');
    const groups = {};
    if (queue) {
      const queueItems: ListItem[] = JSON.parse(queue);
      queueItems.forEach((item: ListItem) => {
        // group by list ID
        if (groups[item.listId]) {
          groups[item.listId].push(item);
        } else {
          groups[item.listId] = [item];
        }
      });
    }

    // eslint-disable-next-line guard-for-in
    for (const listID in groups) {
      const update: ListUpdate = { updates: [], updatedAt: 0 };
      update.updates = [groups[listID]];
      requests$.push(
        this.http.patch(`${environment.api}/lists/${listID}/updates`, update)
      );
    }
    if (requests$.length === 0) {
      this.syncing.next(false);
      return;
    }
    forkJoin(requests$).subscribe(() => {
      localStorage.removeItem('queue');
      this.syncing.next(false);
    });
  }

  enqueueUpdate(item: ListItem) {
    // note: we don't need to coalesce the queue.
    // the server will apply updates one by one, so no issue.
    if (this.syncing.getValue()) {
      console.warn('Enqueuing update when we are already syncing');
    }
    const queue = localStorage.getItem('queue');
    if (!queue) {
      localStorage.setItem('queue', JSON.stringify([item]));
      return;
    }
    const queueItems: ListItem[] = JSON.parse(queue);
    queueItems.push(item);
    localStorage.setItem('queue', JSON.stringify(queueItems));
  }
  ngOnDestroy() {
    if (this.networkListener) {
      this.networkListener.remove();
    }
  }
}
