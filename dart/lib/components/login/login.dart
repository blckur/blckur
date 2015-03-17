library login_comp;

import 'package:blckur/exceptions.dart';
import 'package:blckur/models/models.dart' as models;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

@Component(
  selector: 'x-login',
  templateUrl: 'packages/blckur/components/login/login.html',
  cssUrl: 'packages/blckur/components/login/login.css'
)
class LoginComp implements ng.ShadowRootAware {
  models.Auth model;
  String emailError;
  String passwordError;
  String mode;
  ng.Router router;

  LoginComp(this.model, this.router) : mode = 'login' {
    this.model.remember = true;
  }

  void onShadowRoot(dom.ShadowRoot root) {
    root.querySelector('.modal').onKeyPress.listen((evt) {
      if (evt.keyCode == 13) {
        evt.stopPropagation();
        if (this.mode == 'login') {
          this.onLogin();
        }
        else {
          this.onSignup();
        }
      }
    });
  }

  void _handlerError(HttpError err) {
    this.clearErrors();

    if (err.error == 'email_invalid') {
      this.emailError = 'Email is invalid';
    }
    else if (err.error == 'password_invalid') {
      this.passwordError = 'Password is invalid';
    }
    else if (err.error == 'email_exists') {
      this.emailError = 'Email is already signed up';
    }
  }

  void clearErrors() {
    this.emailError = null;
    this.passwordError = null;
  }

  bool validate() {
    this.clearErrors();

    try {
      this.model.validate('email');
    } catch(err) {
      this.emailError = err.toString();
      return false;
    }

    try {
      this.model.validate('password');
    } catch(err) {
      this.passwordError = err.toString();
      return false;
    }

    return true;
  }

  void setMode(String mode) {
    this.mode = mode;
    this.clearErrors();
  }

  void onLogin() {
    if (!this.validate()) {
      return;
    }

    this.model.login(['email', 'password', 'remember']).then((_) {
      this.clearErrors();
      this.router.gotoUrl('/').then((_) {
        this.model.clear();
      });
    }).catchError((err) {
      logger.severe('Failed to login', err);
      this._handlerError(err);
    });
  }

  void onSignup() {
    if (!this.validate()) {
      return;
    }

    this.model.signup(['email', 'password']).then((_) {
      this.clearErrors();
      this.router.gotoUrl('/root').then((_) {
        this.model.clear();
      });
    }).catchError((err) {
      logger.severe('Failed to signup', err);
      this._handlerError(err);
    });
  }
}
