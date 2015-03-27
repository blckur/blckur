library account_init_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountAdd extends model.Model {
  String type;
  String redirect;

  AccountAdd(ng.Http http) : super(http);

  model.Model newModel() {
    return new AccountAdd(this.http);
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
