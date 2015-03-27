library auth_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

@Injectable()
class Auth extends model.Model {
  String url = '/auth';
  String email;
  void emailValidator(val) {
    if (val == null || val == '') {
      throw new model.Invalid('empty', 'Email cannot be empty');
    }
  }
  String password;
  void passwordValidator(val) {
    if (val == null || val == '') {
      throw new model.Invalid('empty', 'Password cannot be empty');
    }
  }
  bool remember;

  Auth(ng.Http http) : super(http);

  model.Model newModel() {
    return new Auth(this.http);
  }

  Map<String, Function> get getters {
    return {
      'email': () => this.email,
      'password': () => this.password,
      'remember': () => this.remember,
    };
  }

  Map<String, Function> get setters {
    return {
      'email': (x) => this.email = x,
      'password': (x) => this.password = x,
      'remember': (x) => this.remember = x,
    };
  }

  Map<String, Function> get validators {
    return {
      'email': this.emailValidator,
      'password': this.passwordValidator,
    };
  }

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
