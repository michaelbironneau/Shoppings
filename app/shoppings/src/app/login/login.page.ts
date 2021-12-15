import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../shared/services/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
})
export class LoginPage implements OnInit {
  failedLogin = false;
  username: string;
  password: string;
  constructor(private auth: AuthService, private router: Router) {}

  ngOnInit() {}

  onLogin() {
    this.failedLogin = false;
    try {
      console.log(`Logging in as ${this.username}`);
      this.auth
        .doAuth(this.username, this.password)
        .subscribe((success: boolean) => {
          if (!success) {
            console.warn('Failed login');
            this.failedLogin = true;
          } else {
            this.router.navigate(['home']);
          }
        });
    } catch (err) {
      console.error(`Error authenticationg: ${err}`);
      this.failedLogin = true;
    }
  }
}
