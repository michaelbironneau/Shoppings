/* eslint-disable object-shorthand */
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ListItem } from '../shared/models/list-item';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { ListItemService } from '../shared/services/list-item.service';
import { Subscription, timer } from 'rxjs';

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  styleUrls: ['./list.page.scss'],
})
export class ListPage implements OnInit, OnDestroy {
  listID: string;
  items: ListItem[] = [];
  searchString = null;
  searchResults: ListItem[] = [];
  timerSub: Subscription;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private listItemService: ListItemService
  ) {}

  ngOnInit() {
    this.route.params.subscribe((params: Params) => {
      this.listID = params.id;
      console.log(`Editing: ${this.listID}`);
      this.refresh();
      const syncTimer = timer(5000, 10000);
      this.timerSub = syncTimer.subscribe(() => this.asyncUpdate);
    });
  }

  ngOnDestroy() {
    this.timerSub.unsubscribe();
  }

  asyncUpdate() {
    this.listItemService
      .syncListItems(this.listID)
      .subscribe((listItems: ListItem[]) => {
        if (listItems && listItems.length > 0) {
          console.log('Received async update', listItems);
          // Get all from local storage, as sync will update that.
          // What we want to avoid is making a full update request via API.
          this.items = this.listItemService.getAllLocal(this.listID);
        }
      });
  }

  refresh() {
    this.listItemService
      .getAll(this.listID)
      .subscribe((items) => (this.items = items));
  }

  onSearchCancel() {
    this.searchResults = [];
    this.searchString = null;
  }

  onAddCustom() {
    this.listItemService
      .applyUpdate(this.listID, {
        updates: [
          {
            listId: this.listID,
            name: this.searchString,
            quantity: 1,
            checked: false,
          },
        ],
      })
      .subscribe(() => {
        this.refresh();
        this.onSearchCancel();
      });
  }
  onAdd(result) {
    this.listItemService
      .applyUpdate(this.listID, {
        updates: [
          {
            listId: this.listID,
            id: result.id,
            name: result.name,
            quantity: 1,
            checked: false,
          },
        ],
      })
      .subscribe(() => {
        this.refresh();
        this.onSearchCancel();
      });
  }

  onSearch(e) {
    if (e.target.value.length === 0) {
      this.onSearchCancel();
      return;
    }
    this.searchString =
      e.target.value[0].toUpperCase() + e.target.value.substring(1);
    this.listItemService
      .searchAutocomplete(e.target.value)
      .subscribe((results) => {
        this.searchResults = results;
      });
  }

  onShop() {
    this.router.navigate(['/shop', this.listID]);
  }

  onApplyDiff(item: ListItem, quantityDiff: number) {
    const newItem = { ...item };
    newItem.quantity += quantityDiff;
    this.listItemService
      .applyUpdate(this.listID, {
        updates: [newItem],
      })
      .subscribe(() => {
        this.refresh();
      });
  }
}
