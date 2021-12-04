/* eslint-disable object-shorthand */
import { Component, OnInit } from '@angular/core';
import { ListItem } from '../shared/models/list-item';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { ListItemService } from '../shared/services/list-item.service';

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
        name: this.searchString,
        quantityDiff: 1,
      })
      .subscribe(() => {
        this.refresh();
        this.onSearchCancel();
      });
  }
  onAdd(result) {
    this.listItemService
      .applyUpdate(this.listID, {
        name: result.name,
        quantityDiff: 1,
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
    this.listItemService
      .applyUpdate(this.listID, {
        name: item.name,
        quantityDiff: quantityDiff,
      })
      .subscribe(() => {
        this.refresh();
      });
  }
}
