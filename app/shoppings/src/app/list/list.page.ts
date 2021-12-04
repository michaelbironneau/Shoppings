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
