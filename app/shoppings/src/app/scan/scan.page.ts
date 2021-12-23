/* eslint-disable space-before-function-paren */
/* eslint-disable prefer-arrow/prefer-arrow-functions */
/* eslint-disable guard-for-in */
import { Component, OnDestroy, OnInit } from '@angular/core';
import { createWorker, Word } from 'tesseract.js';
import * as foodKeywords from '../shared/data/food-keywords.json';
import * as itemKeywords from '../shared/data/item-keywords.json';
import * as foodItems from '../shared/data/common-foods.json';
import * as itemItems from '../shared/data/items.json';
import { serialize } from 'v8';
import { ListItemService } from '../shared/services/list-item.service';
import { ActivatedRoute, Params, Router } from '@angular/router';

@Component({
  selector: 'app-scan',
  templateUrl: './scan.page.html',
  styleUrls: ['./scan.page.scss'],
})
export class ScanPage implements OnInit, OnDestroy {
  content = null;
  itemCache: Set<string> = new Set();
  allItems: string[] = [];
  progress = 1;
  step = 1;
  results;
  listID: string;
  worker = createWorker({
    logger: (m) => {
      if (m.progress) {
        this.progress = m.progress;
      }
    }, // Add logger here
  });
  constructor(
    private listItemService: ListItemService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    for (const ix in foodKeywords) {
      this.itemCache.add(foodKeywords[ix]);
    }
    for (const ix in itemKeywords) {
      this.itemCache.add(itemKeywords[ix]);
    }
    for (const ix in foodItems) {
      this.allItems.push(foodItems[ix]);
    }
    for (const ix in itemItems) {
      this.allItems.push(itemItems[ix]);
    }
  }

  ngOnInit() {
    this.loadWorker().then(() => {
      console.log('Worker loaded');
    });
    this.route.params.subscribe((params: Params) => {
      this.listID = params.id;
    });
  }

  ngOnDestroy() {
    this.destroyWorker().then(() => {
      console.log('Tesseract worker destroyed');
    });
  }

  async loadWorker() {
    await this.worker.load();
    await this.worker.loadLanguage('eng');
    await this.worker.initialize('eng');
  }

  async destroyWorker() {
    await this.worker.terminate();
  }

  keywordMatches(words: Word[]): string[] {
    const matches = new Set<string>();
    words.forEach((word: Word) => {
      const wordNoPunct = word.text
        .toLowerCase()
        .replace(/[.,\/#!$%\^&\*;:{}=\-_`~()]/g, '');
      if (this.itemCache.has(wordNoPunct)) {
        matches.add(wordNoPunct);
      }
    });
    return Array.from(matches);
  }

  itemMatches(text: string): string[] {
    const matches = new Set<string>();
    this.allItems.forEach((item: string) => {
      const ix = text.indexOf(item);
      if (ix !== -1) {
        matches.add(item);
      }
    });
    return Array.from(matches);
  }

  onImport() {
    console.log(this.results);
    const updates = this.results
      .filter((result) => result.checked)
      .map((result) => ({
        listId: this.listID,
        id: null,
        name: result.name,
        quantity: 1,
        checked: false,
      }));
    this.listItemService
      .applyUpdate(this.listID, {
        updates,
      })
      .subscribe(() => {
        this.router.navigate(['/list', this.listID]);
      });
  }

  toTitleCase(str) {
    return str.replace(/\w\S*/g, function (txt: string) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    });
  }

  onScan(filename: string) {
    console.log('Scanning', filename);
    this.progress = 0;
    this.step = 1;
    this.worker.recognize(`../assets/test-images/${filename}`).then((data) => {
      console.log(data);
      const matches = data.data.lines.map((line) => ({
        text: line.text,
        keywords: this.keywordMatches(line.words),
      }));
      const results = matches
        .filter((match) => match.keywords.length > 0)
        .map((match) => ({
          text: match.text,
          keywords: match.keywords,
          items: this.itemMatches(match.text),
        }));
      const resultSet = new Set<string>(); // to hold the phrases
      const resultKeywords = new Set<string>(); // to hold the phrases split out by words, deduplicating from keywords
      results.forEach((result) => {
        result.items.forEach((item) => {
          resultSet.add(item);
          item.split(' ').forEach((word) => resultKeywords.add(word));
        }); // 1) Add the phrase
        result.keywords.forEach((keyword) => {
          if (!resultKeywords.has(keyword)) {
            resultSet.add(keyword);
            resultKeywords.add(keyword);
          }
        });
      });
      const resultArray = Array.from(resultSet).map((str) =>
        this.toTitleCase(str)
      );
      resultArray.sort();
      this.results = resultArray.map((str) => ({ name: str, checked: true }));
      this.step = 2;
    });
  }
}
