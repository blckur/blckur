library user_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class User extends model.Model {
  String id;
  String email;
  String password;

  User(ng.Http http) : super(http);

  model.Model newModel() {
    return new User(this.http);
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'email': () => this.email,
      'password': () => this.password,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'email': (x) => this.email = x,
      'password': (x) => this.password = x,
    };
  }

  String url = '/user';
}
