library account_types_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

class AccountTypes extends collection.Collection {
  model.Model newModel() {
    return new models.AccountType();
  }

  collection.Collection newCollection() {
    return new AccountTypes();
  }

  String url = '/account_types';
}
