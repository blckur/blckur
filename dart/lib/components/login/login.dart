library login_comp;

import 'package:blckur/models/auth.dart' as ath;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

@Component(
  selector: 'x-login',
  templateUrl: 'packages/blckur/components/login/login.html',
  cssUrl: 'packages/blckur/components/login/login.css'
)
class LoginComp {
  ath.Auth model;
  String emailError;
  String passwordError;
  String mode;
  ng.Router router;

  LoginComp(this.model, this.router) : mode = 'login';

  void _handlerError(async.AsyncError err) {
    if (err.error == 'email_invalid') {
      this.emailError = 'Email is invalid';
      this.passwordError = null;
    }
    else if (err.error == 'password_invalid') {
      this.passwordError = 'Password is invalid';
      this.emailError = null;
    }
    else if (err.error == 'email_exists') {
      this.emailError = 'Email is already signed up';
      this.passwordError = null;
    }
  }

  void onLogin() {
    this.model.login().then((_) {
      this.emailError = null;
      this.passwordError = null;
      this.router.gotoUrl('/root').then((_) {
        this.model.clear();
      });
    }).catchError((err) {
      logger.severe('Failed to login', err);
      this._handlerError(err);
    });
  }

  void onSignup() {
    this.model.signup().then((_) {
      this.emailError = null;
      this.passwordError = null;
      this.router.gotoUrl('/root').then((_) {
        this.model.clear();
      });
    }).catchError((err) {
      logger.severe('Failed to signup', err);
      this._handlerError(err);
    });
  }
}
