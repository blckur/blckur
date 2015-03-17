library account_init_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountAdd extends model.Model {
  String type;

  @model.Attribute('redirect')
  String redirect;

  AccountAdd(ng.Http http) : super(http);

  String get url {
    return '/accounts/${this.type}';
  }
}
