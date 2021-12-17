import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { ToastController } from '@ionic/angular';
import { Subscription, timer } from 'rxjs';
import { ListItem } from '../shared/models/list-item';
import { ListItemService } from '../shared/services/list-item.service';
import { ListService } from '../shared/services/list.service';

@Component({
  selector: 'app-shop',
  templateUrl: './shop.page.html',
  styleUrls: ['./shop.page.scss'],
})
export class ShopPage implements OnInit, OnDestroy {
  listID: string;
  items: ListItem[] = [];
  progress = 0;
  completed = 0;
  total = 0;
  timerSub: Subscription;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private listItemService: ListItemService,
    private listService: ListService,
    private toastController: ToastController
  ) {}

  ngOnInit() {
    this.route.params.subscribe((params: Params) => {
      this.listID = params.id;
      this.refresh();
      const syncTimer = timer(5000, 10000);
      this.timerSub = syncTimer.subscribe(() => this.asyncUpdate());
    });
  }

  ngOnDestroy() {
    this.timerSub.unsubscribe();
  }

  onComplete(item: ListItem) {
    this.listItemService.completeItem(this.listID, item).subscribe(() => {
      this.refresh();
    });
  }

  updateChangesThings(update: ListItem[]): number {
    const yes = update.filter((item: ListItem) => {
      const myIndex = this.items.findIndex(
        (myItem: ListItem) =>
          (myItem.id && item.id && myItem.id === item.id) ||
          myItem.name === item.name
      );
      if (myIndex === -1 && item.quantity === 0) {
        return false; //already deleted
      }
      if (myIndex === -1) {
        return true;
      }
      if (this.items[myIndex].quantity !== item.quantity) {
        return true;
      }
      if (this.items[myIndex].checked !== item.checked) {
        return true;
      }
      return false;
    });
    return yes.length;
  }

  asyncUpdate() {
    this.listItemService
      .syncListItems(this.listID)
      .subscribe((listItems: ListItem[]) => {
        if (listItems && listItems.length > 0) {
          console.log('Received async update', listItems);
          // Get all from local storage, as sync will update that.
          // What we want to avoid is making a full update request via API.
          const changes: number = this.updateChangesThings(listItems);
          this.items = this.listItemService.getAllLocal(this.listID);
          if (changes > 0) {
            this.refreshProgress();
            this.presentUpdateToast(changes);
          }
        }
      });
  }

  async presentUpdateToast(updateLength: number) {
    const toast = await this.toastController.create({
      message: `${updateLength} update${updateLength > 1 ? 's' : ''} received`,
      duration: 2000,
    });
    toast.present();
  }

  refreshProgress() {
    this.total = this.items.length;
    this.completed = 0;
    if (this.total === 0) {
      return; // don't divide by zero later
    }
    this.items.forEach((item) => {
      if (item.checked) {
        this.completed++;
      }
    });
    this.progress = this.completed / this.total;
    // eslint-disable-next-line no-console
    console.debug(`New progress: ${this.progress}`);
  }

  onArchiveList() {
    this.listService.archive(this.listID).subscribe(() => {
      this.router.navigate(['/home']);
    });
  }

  onNavigateList() {
    this.router.navigate(['/list', this.listID]);
  }

  refresh() {
    this.listItemService.getAll(this.listID).subscribe((items) => {
      this.items = items;
      this.refreshProgress();
      // eslint-disable-next-line no-console
      console.debug(`New progress: ${this.progress}`);
    });
  }
}
