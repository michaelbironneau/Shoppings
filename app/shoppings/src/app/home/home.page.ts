/* eslint-disable no-console */
import { Component, OnInit } from '@angular/core';
import { List } from '../shared/models/list';
import { ListService } from '../shared/services/list.service';
import { Router } from '@angular/router';
import { AuthService } from '../shared/services/auth.service';
import { ToastController } from '@ionic/angular';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
})
export class HomePage implements OnInit {
  lists: List[] = [];

  constructor(
    private listService: ListService,
    private router: Router,
    private auth: AuthService,
    private toasts: ToastController
  ) {}

  ngOnInit() {
    this.refresh(null);
  }

  async presentErrorToast(msg: string) {
    const toast = await this.toasts.create({
      header: 'Error',
      message: msg,
      position: 'bottom',
      duration: 5000,
      color: 'danger',
    });
    await toast.present();
  }

  refresh(e) {
    this.listService.getAll().subscribe(
      (lists) => {
        this.lists = lists;
        if (e) {
          e.target.complete();
        }
      },
      (err) => {
        if (err.status === 401) {
          if (e) {
            e.target.complete();
          }
          this.auth.logout();
          this.router.navigate(['/login']);
        } else if (err.status === 0) {
          this.presentErrorToast(
            `Server doesn't seem to be running. Try again shortly and if it still doesn't work, ask Michael to turn it on.`
          );
          if (e) {
            e.target.complete();
          }
        } else if (err.status === 400) {
          this.presentErrorToast(
            `This is an odd error to receive, make sure you're running the latest version of the app.`
          );
          if (e) {
            e.target.complete();
          }
        } else {
          this.presentErrorToast(
            `Please retry soon, and if it still doesn't work tell Michael there was a server error.`
          );
          if (e) {
            e.target.complete();
          }
        }
      }
    );
  }

  onArchive(listID: string) {
    console.debug(`Archiving list ${listID}`);
    this.listService.archive(listID).subscribe((success: boolean) => {
      console.debug(`Success? ${success}`);
      this.refresh(null);
    });
  }

  onOpen(list: List) {
    this.router.navigate(['/list', list.id]);
  }

  uniqueName(): string {
    // come up with a unique name of form 'Dec 12 List #2'
    const startingPoint =
      new Date().toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      }) + ' List';
    let current = startingPoint;
    let i = 1;
    while (i++) {
      const ix = this.lists.findIndex((list) => list.name === current);
      if (ix === -1) {
        break;
      }
      current = startingPoint + ' #' + i.toString();
    }
    return current;
  }

  onAdd() {
    console.debug('Adding new list');
    const newList: List = {
      id: '',
      name: this.uniqueName(),
      archived: false,
      summary: '(empty)',
    };
    this.listService.add(newList).subscribe((newID: string) => {
      console.debug(`New list ID: ${newID}`);
      this.refresh(null);
    });
  }
}
