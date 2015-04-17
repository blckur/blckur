library user_mod;

import 'package:blckur/model.dart' as model;

class User extends model.Model {
  String id;
  String email;

  String password;
  void passwordValidator(val) {
    if (val == null || val == '') {
      throw new model.Invalid('empty', 'Password cannot be empty');
    }
  }

  model.Model newModel() {
    return new User();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'email': () => this.email,
      'password': () => this.password,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'email': (x) => this.email = x,
      'password': (x) => this.password = x,
    };
  }

  Map<String, Function> get validators {
    return {
      'password': this.passwordValidator,
    };
  }

  String url = '/user';
}
