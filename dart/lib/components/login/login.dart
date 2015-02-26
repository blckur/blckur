library login_comp;

import 'package:blckur/models/auth.dart' as ath;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;

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

  LoginComp(this.model) : mode = 'login';

  void onLogin() {
    this.model.login().then((_) {
      this.emailError = null;
      this.passwordError = null;
    }).catchError((err) {
      logger.severe('Failed to login', err);

      if (err.error == 'auth_email_invalid') {
        this.emailError = 'Email is invalid';
        this.passwordError = null;
      }
      else if (err.error == 'auth_password_invalid') {
        this.passwordError = 'Password is invalid';
        this.emailError = null;
      }
    });
  }

  void onSignup() {
    this.model.signup().then((_) {
      this.emailError = null;
      this.passwordError = null;
    }).catchError((err) {
      logger.severe('Failed to signup', err);

      if (err.error == 'auth_email_invalid') {
        this.emailError = 'Email is invalid';
        this.passwordError = null;
      }
      else if (err.error == 'auth_password_invalid') {
        this.passwordError = 'Password is invalid';
        this.emailError = null;
      }
    });
  }
}
