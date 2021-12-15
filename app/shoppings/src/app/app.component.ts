import { Component, OnInit } from '@angular/core';
import { SyncService } from './shared/services/sync.service';

@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  styleUrls: ['app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(private sync: SyncService) {}

  ngOnInit() {
    this.sync.pushQueuedUpdates();
  }
}
