import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { catchError, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { HttpClient } from '@angular/common/http';
import { Token } from '../models/token';
import { SyncService } from './sync.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(private http: HttpClient, private sync: SyncService) {}

  getToken(): string {
    return localStorage.getItem('token');
  }

  logout() {
    localStorage.removeItem('token');
  }

  isAuthenticated(): boolean {
    if (!environment.api) {
      return true;
    }
    const token = this.getToken();
    return token && token.length > 0;
  }

  doAuth(username: string, password: string): Observable<boolean> {
    if (!environment.api) {
      return of(true);
    }
    if (!this.sync.haveNetworkConnectivity.getValue()) {
      throw Error('You must have network connectivity to log in');
    }
    return this.http
      .post<Token>(`${environment.api}/token`, {
        username,
        password,
      })
      .pipe(
        map((token: Token) => {
          localStorage.setItem('token', token.token);
          return true;
        }),
        catchError((err) => {
          console.error(err);
          return of(false);
        })
      );
  }
}
