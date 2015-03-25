library account_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Account extends model.Model {
  @model.Attribute('id')
  String id;

  @model.Attribute('user_id')
  String userId;

  @model.Attribute('type')
  String type;

  @model.Attribute('identity')
  String identity;

  @model.Attribute('alerts')
  List<Map<String, String>> alerts;

  Account(ng.Http http) : super(http);

  String get url {
    return '/accounts/' + this.id;
  }
}
