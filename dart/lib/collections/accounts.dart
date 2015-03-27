library accounts_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Accounts extends collection.Collection {
  Accounts(ng.Http http) : super(http);

  model.Model newModel() {
    return new models.Account(this.http);
  }

  collection.Collection newCollection() {
    return new Accounts(this.http);
  }

  String url = '/accounts';
}
