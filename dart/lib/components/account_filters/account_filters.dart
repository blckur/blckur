library account_filters_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/injectables/injectables.dart' as injectables;
import 'package:blckur/enter_aware.dart' as enter_aware;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgTwoWay;
import 'dart:async' as async;

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
  int _mode;
  async.Future _modeWait;
  bool showFilters;
  bool showAddFilter;
  bool selected;
  dynamic typeValue;
  models.FilterType typeModel;
  collections.FilterTypes filterTypes;
  injectables.Loading loading;

  @NgTwoWay('model')
  models.Account model;

  AccountFiltersComp(this.filterTypes, this.loading) {
    this.mode = FILTERS;
  }

  void set mode (int mode) {
    if (this._modeWait != null) {
      this._modeWait.then((_) {
        this._modeWait = null;
        this.mode = mode;
      });
      return;
    }

    if (this._mode == mode) {
      return;
    }
    this._mode = mode;

    switch (mode) {
      case FILTERS:
        this.showAddFilter = false;
        break;
      case ADD_FILTER:
        this.showFilters = false;
        break;
      default:
        this.showFilters = false;
        this.showAddFilter = false;

        this._modeWait = new async.Future.delayed(
          const Duration(milliseconds: 400), () {
            this._modeWait = null;
          });

        return;
    }

    this._modeWait = new async.Future.delayed(
      const Duration(milliseconds: 800), () {});

    new async.Timer(const Duration(milliseconds: 400), () {
      switch (mode) {
        case FILTERS:
          this.showFilters = true;
          break;
        case ADD_FILTER:
          this.showAddFilter = true;
          break;
      }
      this.mode = mode;
      this._modeWait = null;
    });
  }

  int get mode {
    return this._mode;
  }

  void clear() {
    this.selected = false;
    this.typeValue = null;
  }

  void onAdd() {
    if (!this.loading.set()) {
      return;
    }
    this.mode = NONE;

    this.filterTypes.acctType = this.model.type;
    this.filterTypes.fetch().catchError((err) {
      logger.severe('Failed to load filter types', err);
      new alert.Alert('Failed to load filter types');
    }).whenComplete(() {
      this.mode = ADD_FILTER;
      this.loading.clear();
    });
  }

  void onCancel() {
    this.mode = FILTERS;
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
    this.mode = NONE;

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
      this.mode = FILTERS;
      this.clear();
      this.loading.clear();
    });
  }
}
