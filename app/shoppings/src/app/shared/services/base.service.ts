import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { Network } from '@capacitor/network';
import { from } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class BaseService {
  constructor(protected http: HttpClient) {}

  haveInternet(): Observable<boolean> {
    const status = from(Network.getStatus());
    return status.pipe(map((s) => s.connected));
  }
}
