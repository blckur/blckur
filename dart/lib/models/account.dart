library account_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Account extends mdl.Model {
  @mdl.Attribute('id')
  String id;

  @mdl.Attribute('user_id')
  String userId;

  @mdl.Attribute('type')
  String type;

  @mdl.Attribute('identity')
  String identity;

  @mdl.Attribute('alerts')
  List<Map<String, String>> alerts;

  @mdl.Attribute('resources')
  List<dynamic> resources;

  Account(ng.Http http) : super(http);

  String get url {
    return '/accounts/' + this.id;
  }
}
