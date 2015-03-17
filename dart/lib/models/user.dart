library user_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class User extends model.Model {
  @model.Attribute('id')
  String id;

  @model.Attribute('email')
  String email;

  @model.Attribute('password')
  String password;

  User(ng.Http http) : super(http);

  String url = '/user';
}
