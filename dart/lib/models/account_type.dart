library account_type_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountType extends mdl.Model {
  @mdl.Attribute('label')
  String label;

  @mdl.Attribute('type')
  String type;

  AccountType(ng.Http http) : super(http);
}
