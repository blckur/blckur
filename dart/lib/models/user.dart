library user_mod;

import 'package:blckur/model.dart' as model;

class User extends model.Model {
  String id;
  String email;
  String password;

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

  String url = '/user';
}
