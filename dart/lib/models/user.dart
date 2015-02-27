library user_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class User extends mdl.Model {
  @mdl.Attribute('id')
  String id;

  @mdl.Attribute('email')
  String email;

  @mdl.Attribute('password')
  String password;

  User(ng.Http http) : super(http);

  String url = '/user';
}
