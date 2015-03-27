library account_types_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class AccountTypes extends collection.Collection {
  AccountTypes(ng.Http http) : super(http);

  model.Model newModel() {
    return new models.AccountType(this.http);
  }

  collection.Collection newCollection() {
    return new AccountTypes(this.http);
  }

  String url = '/account_types';
}
