library account_mod;

import 'package:blckur/model.dart' as model;

class Account extends model.Model {
  String id;
  String userId;
  String type;
  String identity;
  bool newAcct;
  List<Map<String, String>> filters;

  model.Model newModel() {
    return new Account();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'user_id': () => this.userId,
      'type': () => this.type,
      'identity': () => this.identity,
      'filters': () => this.filters,
      'new': () => this.newAcct,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'user_id': (x) => this.userId = x,
      'type': (x) => this.type = x,
      'identity': (x) => this.identity = x,
      'filters': (x) => this.filters = x,
      'new': (x) => this.newAcct = x,
    };
  }

  String get url {
    return '/accounts/' + this.id;
  }
}
