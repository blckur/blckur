library account_type_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountType extends model.Model {
  @model.Attribute('id')
  int id;

  @model.Attribute('label')
  String label;

  @model.Attribute('type')
  String type;

  AccountType(ng.Http http) : super(http);
}
