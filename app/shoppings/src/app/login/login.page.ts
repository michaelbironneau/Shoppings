import { Component, OnInit } from '@angular/core';
import { AuthService } from '../shared/services/auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
})
export class LoginPage implements OnInit {
  failedLogin = false;
  constructor(private auth: AuthService) {}

  ngOnInit() {}

  doLogin(username: string, password: string) {
    try {
      if (!this.auth.doAuth(username, password)) {
        this.failedLogin = true;
      }
    } catch (err) {
      console.error(`Error authenticationg: ${err}`);
    }
  }
}
