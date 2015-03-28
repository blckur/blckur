library account_filters_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/enter_aware.dart' as enter_aware;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgTwoWay;

const NONE = 0;
const FILTERS = 1;
const ADD_FILTER = 2;

@Component(
  selector: 'x-account-filters',
  templateUrl:
    'packages/blckur/components/account_filters/account_filters.html',
  cssUrl: 'packages/blckur/components/account_filters/account_filters.css'
)
class AccountFiltersComp extends enter_aware.EnterAware {
  bool selected;
  dynamic typeValue;
  models.FilterType typeModel;
  collections.FilterTypes filterTypes;
  utils.Loading loading;
  utils.ModeSwitch modeswitch;

  @NgTwoWay('model')
  models.Account model;

  AccountFiltersComp() {
    this.filterTypes = new collections.FilterTypes();
    this.loading = new utils.Loading();
    this.modeswitch = new utils.ModeSwitch();
    this.modeswitch.mode = FILTERS;
  }

  void clear() {
    this.selected = false;
    this.typeValue = null;
  }

  void onAdd() {
    if (!this.loading.set()) {
      return;
    }
    this.modeswitch.mode = NONE;

    this.filterTypes.acctType = this.model.type;
    this.filterTypes.fetch().catchError((err) {
      logger.severe('Failed to load filter types', err);
      new alert.Alert('Failed to load filter types');
    }).whenComplete(() {
      this.modeswitch.mode = ADD_FILTER;
      this.loading.clear();
    });
  }

  void onCancel() {
    this.modeswitch.mode = FILTERS;
    this.clear();
    this.loading.clear();
  }

  void onClick(models.FilterType model) {
    this.typeModel = model;

    if (model.valueType == "" || model.valueType == null) {
      this.onSave();
      return;
    }

    this.selected = true;
  }

  void onEnter() {
    if (this.selected) {
      this.onSave();
    }
  }

  void onDel(Map<String, String> alrt) {
    if (!this.loading.set()) {
      return;
    }

    this.model.filters.remove(alrt);
    this.model.save(['filters']).catchError((err) {
      logger.severe('Failed to remove filter', err);
      new alert.Alert('Failed to remove filter');
    }).whenComplete(() {
      this.loading.clear();
    });
  }

  void onSave() {
    if (!this.loading.set()) {
      return;
    }
    this.modeswitch.mode = NONE;

    if (this.model.filters == null) {
      this.model.filters = [];
    }

    this.model.filters.add({
      'type': typeModel.type,
      'value': this.typeValue,
    });

    this.model.save(['filters']).catchError((err) {
      logger.severe('Failed to add filter', err);
      new alert.Alert('Failed to add filter');
    }).whenComplete(() {
      this.modeswitch.mode = FILTERS;
      this.clear();
      this.loading.clear();
    });
  }
}
