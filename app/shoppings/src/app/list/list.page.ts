/* eslint-disable object-shorthand */
import { Component, OnInit } from '@angular/core';
import { ListItem } from '../shared/models/list-item';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { ListItemService } from '../shared/services/list-item.service';
import { Item } from '../shared/models/item';

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  styleUrls: ['./list.page.scss'],
})
export class ListPage implements OnInit {
  listID: string;
  items: ListItem[] = [];
  searchString = null;
  searchResults: ListItem[] = [];
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private listItemService: ListItemService
  ) {}

  ngOnInit() {
    this.route.params.subscribe((params: Params) => {
      this.listID = params.id;
      this.refresh();
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
        console.log(this.searchResults);
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
