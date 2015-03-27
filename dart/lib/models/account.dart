library account_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Account extends model.Model {
  String id;
  String userId;
  String type;
  String identity;
  List<Map<String, String>> alerts;

  Account(ng.Http http) : super(http);

  model.Model newModel() {
    return new Account(this.http);
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'user_id': () => this.userId,
      'type': () => this.type,
      'identity': () => this.identity,
      'alerts': () => this.alerts,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'user_id': (x) => this.userId = x,
      'type': (x) => this.type = x,
      'identity': (x) => this.identity = x,
      'alerts': (x) => this.alerts = x,
    };
  }

  String get url {
    return '/accounts/' + this.id;
  }
}
