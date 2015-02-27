library auth_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

@Injectable()
class Auth extends mdl.Model {
  String url = '/auth';

  @mdl.Attribute('email')
  String email;

  @mdl.Validator('email')
  void emailValidator(val) {
    if (val == null || val == '') {
      throw new mdl.Invalid('empty', 'Email cannot be empty');
    }
  }

  @mdl.Attribute('password')
  String password;

  @mdl.Validator('password')
  void passwordValidator(val) {
    if (val == null || val == '') {
      throw new mdl.Invalid('empty', 'Password cannot be empty');
    }
  }

  @mdl.Attribute('remember')
  bool remember;

  Auth(ng.Http http) : super(http);

  async.Future login([List<String> fields]) {
    this.url = '/login';
    return this.create(fields);
  }

  async.Future logout() {
    this.url = '/session';
    return this.destroy();
  }

  async.Future signup([List<String> fields]) {
    this.url = '/signup';
    return this.create(fields);
  }
}
