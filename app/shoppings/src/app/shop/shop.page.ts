import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { ListItem } from '../shared/models/list-item';
import { ListItemService } from '../shared/services/list-item.service';
import { ListService } from '../shared/services/list.service';

@Component({
  selector: 'app-shop',
  templateUrl: './shop.page.html',
  styleUrls: ['./shop.page.scss'],
})
export class ShopPage implements OnInit {
  listID: string;
  items: ListItem[] = [];
  progress = 0;
  completed = 0;
  total = 0;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private listItemService: ListItemService,
    private listService: ListService
  ) {}

  ngOnInit() {
    this.route.params.subscribe((params: Params) => {
      this.listID = params.id;
      this.refresh();
    });
  }

  onComplete(item: ListItem) {
    this.listItemService.completeItem(this.listID, item).subscribe(() => {
      this.refresh();
    });
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
    });
  }
}
