library account_type_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountType extends model.Model {
  int id;
  String label;
  String type;

  AccountType(ng.Http http) : super(http);

  model.Model newModel() {
    return new AccountType(this.http);
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'label': () => this.label,
      'type': () => this.type,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'label': (x) => this.label = x,
      'type': (x) => this.type = x,
    };
  }
}
