import { Observable, of } from 'rxjs';
import { Injectable } from '@angular/core';
import { BaseService } from './base.service';
import { catchError, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment';
import { HttpClient } from '@angular/common/http';
import { Token } from '@angular/compiler';

@Injectable({
  providedIn: 'root',
})
export class AuthService extends BaseService {
  constructor(private http: HttpClient) {
    super();
  }

  getToken(): string {
    return localStorage.getItem('token');
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
    if (!this.haveNetworkConnectivity) {
      throw Error('You must have network connectivity to log in');
    }
    return this.http
      .post<Token>(`${environment.api}/token`, {
        username,
        password,
      })
      .pipe(
        map((token: Token) => {
          localStorage.set('token', token);
          return true;
        }),
        catchError((err) => of(false))
      );
  }
}
