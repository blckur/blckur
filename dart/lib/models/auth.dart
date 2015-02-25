library auth_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Auth extends mdl.Model {
  @mdl.Attribute('email')
  String email;

  @mdl.Validator('email')
  void emailValidator(val) {
    if (val == null || val == '') {
      throw new mdl.Invalid('empty', 'Email cannot be empty');
    }
  }

  @mdl.Attribute('password')
  String password;

  @mdl.Validator('password')
  void passwordValidator(val) {
    if (val == null || val == '') {
      throw new mdl.Invalid('empty', 'Password cannot be empty');
    }
  }

  Auth(ng.Http http) : super(http);

  String url = '/login';
}
