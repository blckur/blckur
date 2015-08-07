library user_mod;

import 'package:blckur/model.dart' as model;

import 'dart:async' as async;

class User extends model.Model {
  String id;

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

  String apikey;

  model.Model newModel() {
    return new User();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'email': () => this.email,
      'password': () => this.password,
      'apikey': () => this.apikey,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'email': (x) => this.email = x,
      'password': (x) => this.password = x,
      'apikey': (x) => this.apikey = x,
    };
  }

  Map<String, Function> get validators {
    return {
      'email': this.emailValidator,
      'password': this.passwordValidator,
    };
  }

  String url = '/user';

  async.Future refreshApi() {
    this.url = '/user/apikey';
    return this.save();
  }
}
