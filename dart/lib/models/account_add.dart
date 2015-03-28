library account_init_mod;

import 'package:blckur/model.dart' as model;

class AccountAdd extends model.Model {
  String type;
  String redirect;

  model.Model newModel() {
    return new AccountAdd();
  }

  Map<String, Function> get getters {
    return {
      'redirect': () => this.redirect,
    };
  }

  Map<String, Function> get setters {
    return {
      'redirect': (x) => this.redirect = x,
    };
  }

  String get url {
    return '/accounts/${this.type}';
  }
}
