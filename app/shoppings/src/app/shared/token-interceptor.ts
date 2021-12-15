import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(private router: Router) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    // if the request is not to an Open Energi website, skip
    if (!environment.api) {
      console.warn('Trying to make HTTP requests in local seed mode');
      return next.handle(req);
    }
    if (req.url.indexOf('token') !== -1) {
      return next.handle(req);
    }

    const token = localStorage.getItem('token');

    // if there is no token, skip
    if (!token) {
      console.warn('No token - redirecting to login');
      this.router.navigate(['/login']);
      return of(null);
    }

    // Add token to headers
    const modified = req.clone({
      headers: req.headers.set('X-Token', token),
    });

    return next.handle(modified);
  }
}
