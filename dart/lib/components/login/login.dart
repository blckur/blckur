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

  LoginComp(this.model);

  void onLogin() {
    this.model.create().then((_) {
    }).catchError((err) {
      logger.severe('Failed to load accounts', err);

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
