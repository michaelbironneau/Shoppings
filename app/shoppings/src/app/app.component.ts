import { Component, OnInit } from '@angular/core';
import { ListService } from './shared/services/list.service';
import { SyncService } from './shared/services/sync.service';

@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  styleUrls: ['app.component.scss'],
})
export class AppComponent implements OnInit {
  lastListSync = 0;
  constructor(private sync: SyncService, private lists: ListService) {}

  ngOnInit() {
    this.sync.pushQueuedUpdates();
    this.sync.haveNetworkConnectivity.subscribe((status: boolean) => {
      console.log('App component syncinc lists...');
      const t = new Date().getTime();
      if (status) {
        this.sync.pushQueuedUpdates();
      }
      if (status && 0.001 * (t - this.lastListSync) > 5) {
        this.lists.trySyncLists().subscribe((success: boolean) => {
          if (success) {
            console.log('Sync successful');
            this.lastListSync = t;
          } else {
            console.log('Sync failed');
          }
        });
      }
    });
  }
}
