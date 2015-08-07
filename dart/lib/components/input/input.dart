library input_comp;

import 'package:blckur/all_aware.dart' as all_aware;

import 'package:angular/angular.dart' show Component, NgAttr, NgTwoWay;
import 'dart:html' as dom;

@Component(
  selector: 'x-input',
  templateUrl: 'packages/blckur/components/input/input.html',
  cssUrl: 'packages/blckur/components/input/input.css'
)
class InputComp extends all_aware.AllAware {
  String label;
  String labelClass;
  bool labelFloat;

  bool isReadonly;
  @NgAttr('readonly')
  void set readonly(String val) {
    if (val == '') {
      this.isReadonly = true;
    } else {
      this.isReadonly = false;
    }
  }

  @NgAttr('type')
  String type;

  @NgAttr('theme')
  String theme;

  var _placeholder;
  @NgAttr('placeholder')
  void set placeholder(String val) {
    this._placeholder = val;

    if (this.error == null) {
      this.label = val;
    }
  }
  String get placeholder {
    return this._placeholder;
  }

  @NgTwoWay('model')
  String model;

  var _error;
  @NgTwoWay('error')
  void set error(String val) {
    if (this._error != val) {
      this._error = val;
      if (val == null) {
        this.label = this.placeholder;
        this.labelClass = null;
        this.labelFloat = false;
      }
      else {
        this.label = val;
        this.labelClass = 'invalid';
        if (this.model != '' && this.model != null) {
          this.labelFloat = true;
        }
      }
    }
  }
  String get error {
    return this._error;
  }

  InputComp() : labelFloat = false;

  void onAll(dom.ShadowRoot root) {
    root.querySelector('input')
      ..onInput.listen((_) {
        this.error = null;
      })
      ..attributes["type"] = this.type;
  }
}
