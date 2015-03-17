library auth_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

@Injectable()
class Auth extends model.Model {
  String url = '/auth';

  @model.Attribute('email')
  String email;

  @model.Validator('email')
  void emailValidator(val) {
    if (val == null || val == '') {
      throw new model.Invalid('empty', 'Email cannot be empty');
    }
  }

  @model.Attribute('password')
  String password;

  @model.Validator('password')
  void passwordValidator(val) {
    if (val == null || val == '') {
      throw new model.Invalid('empty', 'Password cannot be empty');
    }
  }

  @model.Attribute('remember')
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
