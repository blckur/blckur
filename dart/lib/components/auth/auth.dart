library auth_comp;

import 'package:blckur/exceptions.dart';
import 'package:blckur/models/models.dart' as models;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

@Component(
  selector: 'x-auth',
  templateUrl: 'packages/blckur/components/auth/auth.html',
  cssUrl: 'packages/blckur/components/auth/auth.css'
)
class AuthComp implements ng.ShadowRootAware {
  models.Auth model;
  String emailError;
  String passwordError;
  String mode;
  ng.Router router;

  AuthComp(this.router) {
    this.model = new models.Auth();
    this.model.remember = true;
    this.mode = this.router.activePath[0].name;
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

  void _handlerError(HttpError err, String def) {
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
    else {
      new alert.Alert(def);
    }
  }

  void clearErrors() {
    this.emailError = null;
    this.passwordError = null;
  }

  bool validate([bool emailOnly]) {
    this.clearErrors();

    try {
      this.model.validate('email');
    } catch(err) {
      this.emailError = err.toString();
      return false;
    }

    if (emailOnly != true) {
      try {
        this.model.validate('password');
      } catch(err) {
        this.passwordError = err.toString();
        return false;
      }
    }

    return true;
  }

  void setMode(String mode) {
    dom.window.history.pushState(null, "Blckur", '#/' + mode);
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
      this._handlerError(err, 'Error logging in');
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
      this._handlerError(err, 'Error signing up');
    });
  }

  void onForgot() {
    if (!this.validate(true)) {
      return;
    }

    this.model.reset(['email']).then((_) {
      this.setMode('login');
      new alert.Alert('Password reset instructions have been emailed');
    }).catchError((err) {
      logger.severe('Failed to reset password', err);
      this._handlerError(err, 'Error reseting password');
    });
  }

  void onLogout() {
    this.model.logout().then((_) {
      this.setMode('login');
    }).catchError((err) {
      logger.severe('Failed to logout', err);
    });
  }
}
