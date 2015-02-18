library test_mod;

import 'package:blckur/model.dart' as mdl;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Test extends mdl.Model {
  @mdl.Attribute('id')
  String id;

  Test(ng.Http http) : super(http);
}
